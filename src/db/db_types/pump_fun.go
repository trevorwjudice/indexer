package db_types

import "indexer/src/indexer/keycache"

type PumpFunCreate struct {
	Mint                   PublicKey `db:"mint"`
	Deployer               PublicKey `db:"deployer"`
	BondingCurve           PublicKey `db:"bonding_curve"`
	AssociatedBondingCurve PublicKey `db:"associated_bonding_curve"`
	MetadataSlot           PublicKey `db:"metadata_slot"`
	Name                   string    `db:"name"`
	Symbol                 string    `db:"symbol"`
	MetadataURI            string    `db:"metadata_uri"`
	*InstructionMetadata
}

func (c *PumpFunCreate) Table() string {
	return PUMP_FUN_CREATE
}

func (c *PumpFunCreate) Filter(k *keycache.Keycache) bool {
	k.Add(c.Mint.PublicKey())
	return true
}

type PumpFunSwap struct {
	Mint              PublicKey `db:"mint"`
	MakerTokenAccount PublicKey `db:"maker_token_account"`
	Maker             PublicKey `db:"maker"`
	TokenAmount       int64     `db:"token_amount"`
	SolAmount         int64     `db:"sol_amount"`
	Fee               uint64    `db:"fee"`
	*InstructionMetadata
}

func (s *PumpFunSwap) Table() string {
	return PUMP_FUN_SWAPS
}

func (c *PumpFunSwap) Filter(k *keycache.Keycache) bool {
	return true
}

type PumpFunSetParams struct {
	FeeRecipient                PublicKey `db:"fee_recipient"`
	InitialVirtualTokenReserves uint64    `db:"initial_virtual_token_reserves"`
	InitialVirtualSolReserves   uint64    `db:"initial_virtual_sol_reserves"`
	InitialRealTokenReserves    uint64    `db:"initial_real_token_reserves"`
	TokenTotalSupply            uint64    `db:"token_total_supply"`
	FeeBasisPoints              uint64    `db:"fee_basis_points"`
	*InstructionMetadata
}

func (s *PumpFunSetParams) Table() string {
	return PUMP_FUN_SET_PARAMS
}

func (c *PumpFunSetParams) Filter(k *keycache.Keycache) bool {
	return true
}

type PumpFunWithdraw struct {
	*InstructionMetadata
}

func (s *PumpFunWithdraw) Table() string {
	return PUMP_FUN_WITHDRAW
}

func (c *PumpFunWithdraw) Filter(k *keycache.Keycache) bool {
	return true
}
