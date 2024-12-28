package indexer

import (
	"time"

	"github.com/gagliardetto/solana-go"
)

type ParsedInstruction interface {
	Table() string
	GetSlot() uint64
	GetTransactionIndex() int
	GetInstructionIndex() uint8
	GetTimestamp() time.Time
	GetSignature() solana.Signature
}

type InstructionMetadata struct {
	Slot             uint64           `db:"slot"`
	TransactionIndex int              `db:"transaction_index"`
	InstructionIndex uint8            `db:"instruction_index"`
	Timestamp        time.Time        `db:"ts"`
	Signature        solana.Signature `db:"signature"`
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
	return i.Signature
}
