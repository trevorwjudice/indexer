package db_types

import (
	"indexer/src/indexer/keycache"
	"indexer/src/util/solana/transactions"
	"time"

	"github.com/gagliardetto/solana-go"
)

type ParsedInstruction interface {
	Filter(k *keycache.Keycache) bool
	Table() string
	GetSlot() uint64
	GetTransactionIndex() int
	GetInstructionIndex() uint8
	GetTimestamp() time.Time
	GetSignature() solana.Signature
}

type InstructionMetadata struct {
	Slot             uint64    `db:"slot"`
	TransactionIndex int       `db:"transaction_index"`
	InstructionIndex uint8     `db:"instruction_index"`
	Timestamp        time.Time `db:"ts"`
	Signature        Signature `db:"signature"`
}

func (i *InstructionMetadata) GetSlot() uint64 {
	return i.Slot
}

func (i *InstructionMetadata) GetTransactionIndex() int {
	return i.TransactionIndex
}

func (i *InstructionMetadata) GetInstructionIndex() uint8 {
	return i.InstructionIndex
}

func (i *InstructionMetadata) GetTimestamp() time.Time {
	return i.Timestamp
}

func (i *InstructionMetadata) GetSignature() solana.Signature {
	return i.Signature.Signature()
}

func PopulateMetadata(r *transactions.Reader, flatIndex uint8) *InstructionMetadata {
	return &InstructionMetadata{
		Slot:             r.GetSlot(),
		TransactionIndex: r.GetTransactionIndex(),
		InstructionIndex: flatIndex,
		Signature:        Signature(r.GetSignature()),
		Timestamp:        r.GetTimestamp(),
	}
}
