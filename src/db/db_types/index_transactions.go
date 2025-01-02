package db_types

import (
	"time"
)

type Progress struct {
	SlotStart  uint64 `db:"slot_start"`
	SlotEnd    uint64 `db:"slot_end"`
	Status     int64  `db:"status"`
	BlockCount int64  `db:"block_count"`
	TimeTaken  int64  `db:"time_taken"`
}

type PumpfunCreate struct {
	Block                  uint64    `db:"block"`
	TransactionIndex       int       `db:"transaction_index"`
	InstructionIndex       uint8     `db:"instruction_index"`
	Timestamp              time.Time `db:"ts"`
	Signature              string    `db:"signature"`
	Mint                   string    `db:"mint"`
	Deployer               string    `db:"deployer"`
	BondingCurve           string    `db:"bonding_curve"`
	AssociatedBondingCurve string    `db:"associated_bonding_curve"`
	MetadataSlot           string    `db:"metadata_slot"`
	Name                   string    `db:"name"`
	Symbol                 string    `db:"symbol"`
	MetadataURI            string    `db:"metadata_uri"`
}
