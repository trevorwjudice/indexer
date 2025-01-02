package parser

import (
	"context"
	"indexer/src/db/db_types"
	"indexer/src/indexer/keycache"
	"indexer/src/util/solana/transactions"

	"github.com/gagliardetto/solana-go"
	"github.com/upper/db/v4"
)

type ProgramParseFunc func(r *transactions.Reader, flatInstructionIndex uint8, prev []db_types.ParsedInstruction) ([]db_types.ParsedInstruction, error)

type InstructionParser struct {
	programParsers map[solana.PublicKey]ProgramParseFunc
	kc             *keycache.Keycache
}

func New(fetch keycache.FetchFunc) *InstructionParser {
	return &InstructionParser{kc: keycache.New(fetch), programParsers: make(map[solana.PublicKey]ProgramParseFunc)}
}

func (i *InstructionParser) SetProgramParseFunc(program solana.PublicKey, parseFn ProgramParseFunc) {
	i.programParsers[program] = parseFn
}

func (i *InstructionParser) NewSession(ctx context.Context, s db.Session) (*Session, error) {
	err := i.kc.Fetch(ctx, s)
	if err != nil {
		return nil, err
	}
	return &Session{
		parser: i,
	}, nil
}

func (i *InstructionParser) GetProgramParseFn(id solana.PublicKey) (ProgramParseFunc, bool) {
	fn, ok := i.programParsers[id]
	return fn, ok
}
