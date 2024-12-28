package parser

import (
	"encoding/json"

	"github.com/gagliardetto/solana-go/rpc"
)

type Block struct {
	Slot  uint64 `db:"slot"`
	Block []byte `db:"block"`
}

func NewBlock(block rpc.GetBlockResult, slot uint64) (*Block, error) {
	jzn, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}
	return &Block{
		Slot:  slot,
		Block: jzn,
	}, nil
}
