package parser

import (
	"context"
	"indexer/src/db/db_types"
	"indexer/src/util/dbutil"

	"github.com/upper/db/v4"
)

type Inserter interface {
	Execute(ctx context.Context, s db.Session) error
	Add(inst db_types.ParsedInstruction)
}

type ProgramInstructionSorter struct {
	instTable map[string][]db_types.ParsedInstruction
}

func NewProgramInstructionInserter() Inserter {
	return &ProgramInstructionSorter{instTable: make(map[string][]db_types.ParsedInstruction)}
}

func (p *ProgramInstructionSorter) Execute(ctx context.Context, s db.Session) error {
	for table, insts := range p.instTable {
		err := dbutil.BatchUpload(s, table, insts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ProgramInstructionSorter) Add(inst db_types.ParsedInstruction) {
	p.instTable[inst.Table()] = append(p.instTable[inst.Table()], inst)
}
