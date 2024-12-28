package indexer

import (
	"context"
	"indexer/src/util/dbutil"

	"github.com/upper/db/v4"
)

type Inserter interface {
	Execute(ctx context.Context, s db.Session) error
	Add(inst ParsedInstruction)
}

type ProgramInstructionSorter struct {
	instTable map[string][]ParsedInstruction
}

func NewProgramInstructionInserter() Inserter {
	return &ProgramInstructionSorter{instTable: make(map[string][]ParsedInstruction)}
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

func (p *ProgramInstructionSorter) Add(inst ParsedInstruction) {
	p.instTable[inst.Table()] = append(p.instTable[inst.Table()], inst)
}
