package indexer

import (
	"context"
	"errors"
	"fmt"
	"indexer/src/db/db_types"
	"indexer/src/indexer/keycache"
	"indexer/src/indexer/parser"
	"indexer/src/util/dbutil"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/upper/db/v4"
	"golang.org/x/sync/errgroup"
	"tuxpa.in/a/zlog/log"
)

var SLOT_GROUP_SIZE uint64 = 1000

type Indexer struct {
	cl *rpc.Client
	s  db.Session
	p  *parser.InstructionParser
}

func NewIndexer(cl *rpc.Client, s db.Session, fetch keycache.FetchFunc) *Indexer {
	return &Indexer{cl: cl, s: s, p: parser.New(fetch)}
}

func (i *Indexer) AddParser(id solana.PublicKey, fn parser.ProgramParseFunc) {
	i.p.SetProgramParseFunc(id, fn)
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

	s, err := i.p.NewSession(ctx, i.s)
	if err != nil {
		return err
	}

	wg := errgroup.Group{}
	wg.SetLimit(250)
	failCount := 0
	for _, slot := range blocks {
		slot := slot
		wg.Go(func() (err error) {
			for {
				blk, err := i.cl.GetBlockWithOpts(ctx, slot, &rpc.GetBlockOpts{
					MaxSupportedTransactionVersion: rpc.NewTransactionVersion(0),
					Commitment:                     rpc.CommitmentFinalized,
				})
				if err != nil {
					failCount += 1
					fmt.Printf("error getting slot %d: %s\n", slot, err.Error())
					continue
				}
				err = s.ParseSlot(blk, slot)
				if err != nil {
					panic(err)
					return err
				}
				break
			}
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		return err
	}

	if failCount > 1000 {
		panic("fail count exceeded")
	}

	return i.s.TxContext(ctx, func(d db.Session) error {
		err := s.Execute(ctx, d)
		if err != nil {
			return err
		}
		t.BlockCount = int64(len(blocks))
		t.Status = 1
		t.TimeTaken = int64(time.Since(ts).Seconds())
		return i.finishTask(d, t)
	}, nil)
}

var ErrorZeroSlotsBetween error = errors.New("valid slot length is zero")

const MAX_RETRY_COUNT = 100

func (i *Indexer) getValidBlocksBetween(ctx context.Context, slotStart, slotEnd uint64) ([]uint64, error) {
	retryCount := 1
	validBlocks, err := i.cl.GetBlocks(ctx, slotStart, &slotEnd, rpc.CommitmentFinalized)
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
