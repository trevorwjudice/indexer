package parser

import (
	"context"
	"fmt"
	"indexer/src/db/db_types"
	"indexer/src/util/dbutil"
	"indexer/src/util/solana/transactions"
	"sort"
	"sync"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/upper/db/v4"
)

type Session struct {
	m      sync.RWMutex
	parsed []db_types.ParsedInstruction
	parser *InstructionParser
}

func (s *Session) ParseSlot(blk *rpc.GetBlockResult, slot uint64) error {
	for i := range blk.Transactions {
		err := s.ParseTransaction(blk, slot, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Session) ParseTransaction(blk *rpc.GetBlockResult, slot uint64, txIndex int) error {
	if blk.Transactions[txIndex].Meta.Err != nil {
		return nil
	}

	reader, err := transactions.NewReader(blk, slot, txIndex)
	if err != nil {
		return err
	}

	instructions := reader.GetFlattenedInstructions()
	parsedInstructions := make([]db_types.ParsedInstruction, 0, len(instructions))

	for i, inst := range instructions {
		id, err := reader.GetAccountAtIndex(inst.ProgramIDIndex)
		if err != nil {
			return err
		}

		fn, ok := s.parser.GetProgramParseFn(id)
		if !ok {
			continue
		}

		items, err := fn(reader, uint8(i), parsedInstructions)
		if err != nil {
			return err
		}

		if items == nil {
			continue
		}

		parsedInstructions = append(parsedInstructions, items...)
	}

	s.m.Lock()
	defer s.m.Unlock()
	s.parsed = append(s.parsed, parsedInstructions...)

	return nil
}

func (s *Session) Execute(ctx context.Context, d db.Session) error {
	instructions := sortParsedInstructions(s.parsed)
	insertMap := make(map[string][]db_types.ParsedInstruction)
	for _, inst := range instructions {
		if !inst.Filter(s.parser.kc) {
			continue
		}
		insertMap[inst.Table()] = append(insertMap[inst.Table()], inst)
	}

	for table, rows := range insertMap {
		err := dbutil.BatchUpload(d, table, rows)
		if err != nil {
			return err
		}

		fmt.Println("inserted ", len(rows), " rows into table ", table)
	}

	return nil
}

func sortParsedInstructions(instructions []db_types.ParsedInstruction) []db_types.ParsedInstruction {
	sort.SliceStable(instructions, func(i, j int) bool {
		if instructions[i].GetSlot() < instructions[j].GetSlot() {
			return true
		}
		if instructions[i].GetTransactionIndex() < instructions[j].GetTransactionIndex() {
			return true
		}
		return instructions[i].GetInstructionIndex() < instructions[j].GetInstructionIndex()
	})
	return instructions
}
