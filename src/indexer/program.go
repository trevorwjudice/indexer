package indexer

import "indexer/src/util/solana/transactions"

type InstructionParser interface {
	ParseTransaction(r *transactions.Reader) ([]ParsedInstruction, error)
}
