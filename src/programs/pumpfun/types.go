package pumpfun

import (
	"indexer/src/db_types"
	"indexer/src/indexer"

	"github.com/gagliardetto/solana-go"
)

type Create struct {
	Mint                   solana.PublicKey `db:"mint"`
	Deployer               solana.PublicKey `db:"deployer"`
	BondingCurve           solana.PublicKey `db:"bonding_curve"`
	AssociatedBondingCurve solana.PublicKey `db:"associated_bonding_curve"`
	MetadataSlot           solana.PublicKey `db:"metadata_slot"`
	Name                   string           `db:"name"`
	Symbol                 string           `db:"symbol"`
	MetadataURI            string           `db:"metadata_uri"`
	*indexer.InstructionMetadata
}

func (c *Create) Table() string {
	return db_types.PUMP_FUN_CREATE
}

type Swap struct {
	Mint              solana.PublicKey `db:"mint"`
	MakerTokenAccount solana.PublicKey `db:"maker_token_account"`
	Maker             solana.PublicKey `db:"maker"`
	TokenAmount       int64            `db:"token_amount"`
	SolAmount         int64            `db:"sol_amount"`
	Fee               uint64           `db:"fee"`
	*indexer.InstructionMetadata
}

func (s *Swap) Table() string {
	return db_types.PUMP_FUN_SWAPS
}
