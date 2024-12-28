package indexer

import (
	"context"
	"errors"
	"fmt"
	"indexer/src/db_types"
	"indexer/src/util/dbutil"
	"indexer/src/util/solana/transactions"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/upper/db/v4"
	"golang.org/x/sync/errgroup"
	"tuxpa.in/a/zlog/log"
)

var SLOT_GROUP_SIZE uint64 = 1000

type Indexer struct {
	cl      *rpc.Client
	s       db.Session
	parsers []InstructionParser
}

func NewIndexer(cl *rpc.Client, s db.Session) *Indexer {
	return &Indexer{cl: cl, s: s}
}

func (i *Indexer) AddProgramParser(parser InstructionParser) {
	i.parsers = append(i.parsers, parser)
}

func (i *Indexer) GetBackendBlockHeight(ctx context.Context) (uint64, error) {
	return 0, nil
}

func (i *Indexer) Run(ctx context.Context) error {
	// First, reset the status of all "in progress" tasks.
	err := i.resetTaskProgress(ctx)
	if err != nil {
		return err
	}

	for {
		startTs := time.Now()
		task, err := i.startNextTask(ctx)
		if err != nil {
			return err
		}

		log.Info().Uint64("slot_start", task.SlotStart).Uint64("slot_end", task.SlotEnd).Msg("starting task")

		err = i.doTask(ctx, task)
		if err != nil {
			return err
		}

		log.Info().Uint64("slot_start", task.SlotStart).Uint64("slot_end", task.SlotEnd).Dur("duration", time.Since(startTs)).Msg("task executed successfully")
	}
}

func (i *Indexer) ScheduleRange(ctx context.Context, startSlot, endSlot uint64) error {
	if endSlot < startSlot {
		return fmt.Errorf("end slot cannot be greater than starts lot")
	}
	currScheduleHead, err := i.getScheduledSlotHeight(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if currScheduleHead > endSlot {
		return nil
	}

	if currScheduleHead >= startSlot {
		startSlot = currScheduleHead + 1
	}

	tasks := make([]*db_types.Progress, 0, (endSlot-startSlot)/SLOT_GROUP_SIZE+1)
	for i := startSlot; i < endSlot; i += (SLOT_GROUP_SIZE + 1) {
		newTask := &db_types.Progress{
			SlotStart:  i,
			SlotEnd:    i + SLOT_GROUP_SIZE,
			Status:     0,
			TimeTaken:  0,
			BlockCount: 0,
		}
		tasks = append(tasks, newTask)
	}

	err = dbutil.BatchUpload(i.s, db_types.INDEXER_PROGRESS, tasks)
	if err != nil {
		return err
	}

	return nil
}

func (i *Indexer) doTask(ctx context.Context, t *db_types.Progress) error {

	ts := time.Now()
	blocks, err := i.getValidBlocksBetween(ctx, t.SlotStart, t.SlotEnd)
	if err != nil {
		return err
	}

	wg := errgroup.Group{}
	wg.SetLimit(250)
	results := make([]*BlockResult, len(blocks))
	for ind, slot := range blocks {
		ind := ind
		slot := slot
		wg.Go(func() (err error) {
			ts := time.Now()
			for {
				results[ind], err = i.GetSlot(ctx, slot)
				if err != nil {
					fmt.Printf("error getting slot %d: %s\n", slot, err.Error())
					continue
				}
				break
			}
			fmt.Println(time.Since(ts))
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	// upload blocks

	inserter := NewProgramInstructionInserter()

	for _, blk := range results {
		for _, inst := range blk.ParsedInstructions {
			inserter.Add(inst)
		}
	}

	return i.s.TxContext(ctx, func(s db.Session) error {
		err := inserter.Execute(ctx, s)
		if err != nil {
			return err
		}

		t.BlockCount = int64(len(blocks))
		t.Status = 1
		t.TimeTaken = int64(time.Since(ts).Seconds())
		return i.finishTask(s, t)
	}, nil)
}

type BlockResult struct {
	ParsedInstructions []ParsedInstruction
}

type TransactionInfo struct {
	Tx       *rpc.TransactionWithMeta
	Parsed   *solana.Transaction
	Accounts []solana.PublicKey
}

func (i *Indexer) GetSlot(ctx context.Context, slot uint64) (*BlockResult, error) {
	u := uint64(0)
	blk, err := i.cl.GetBlockWithOpts(ctx, slot, &rpc.GetBlockOpts{
		MaxSupportedTransactionVersion: &u,
	})
	if err != nil {
		return nil, err
	}

	res := &BlockResult{}
	for ind, tx := range blk.Transactions {
		if tx.Meta.Err != nil {
			// Possibly record failed TX in the future? Might be useful
			continue
		}

		parsedInstructions, err := i.ParseTransaction(blk, slot, ind)
		if err != nil {
			return nil, err
		}

		for _, inst := range parsedInstructions {
			if inst != nil {
				res.ParsedInstructions = append(res.ParsedInstructions, inst)

			}
		}
	}

	return res, nil
}

func (i *Indexer) ParseTransaction(blk *rpc.GetBlockResult, slot uint64, txIndex int) ([]ParsedInstruction, error) {
	var res []ParsedInstruction

	for _, parser := range i.parsers {
		reader, err := transactions.NewReader(blk, slot, txIndex)
		if err != nil {
			return nil, err
		}
		parsedInstructions, err := parser.ParseTransaction(reader)
		if err != nil {
			return nil, err
		}

		res = append(res, parsedInstructions...)
	}

	return res, nil
}

var ErrorZeroSlotsBetween error = errors.New("valid slot length is zero")

const MAX_RETRY_COUNT = 100

func (i *Indexer) getValidBlocksBetween(ctx context.Context, slotStart, slotEnd uint64) ([]uint64, error) {
	retryCount := 1
	validBlocks, err := i.cl.GetBlocks(ctx, slotStart, &slotEnd, rpc.CommitmentConfirmed)
	for err != nil {
		log.Err(err).Uint64("slot_start", slotStart).Uint64("slot_end", slotEnd).Msg("error getting blocks")
		if retryCount == MAX_RETRY_COUNT {
			return nil, fmt.Errorf("retry count exceeded")
		}
		validBlocks, err = i.cl.GetBlocks(ctx, slotStart, &slotEnd, rpc.CommitmentConfirmed)
		if len(validBlocks) == 0 {
			err = ErrorZeroSlotsBetween
		}
		retryCount += 1
	}
	return []uint64(validBlocks), nil
}
