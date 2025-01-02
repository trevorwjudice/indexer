package db_types

import "indexer/src/indexer/keycache"

type SplTransfer struct {
	Mint        PublicKey `db:"mint"`
	Authority   PublicKey `db:"authority"`
	Source      PublicKey `db:"source"`
	Destination PublicKey `db:"destination"`
	Amount      uint64    `db:"amount"`
	*InstructionMetadata
}

func (s *SplTransfer) Table() string {
	return SPL_TRANSFER
}

func (s *SplTransfer) Filter(k *keycache.Keycache) bool {
	return k.Contains(s.Mint.PublicKey())
}

type SplInitializeAccount struct {
	Owner   PublicKey `db:"owner"`
	Mint    PublicKey `db:"mint"`
	Account PublicKey `db:"account"`
	*InstructionMetadata
}

func (s *SplInitializeAccount) Table() string {
	return SPL_INITIALIZE_ACCOUNT
}

func (s *SplInitializeAccount) Filter(k *keycache.Keycache) bool {
	return k.Contains(s.Mint.PublicKey())
}

type SplBurn struct {
	Mint    PublicKey `db:"mint"`
	Account PublicKey `db:"account"`
	Owner   PublicKey `db:"owner"`
	Amount  uint64    `db:"amount"`
	*InstructionMetadata
}

func (s *SplBurn) Table() string {
	return SPL_BURN
}

func (s *SplBurn) Filter(k *keycache.Keycache) bool {
	return k.Contains(s.Mint.PublicKey())
}

type AssociatedTokenAccountCreate struct {
	Account PublicKey `db:"account"`
	Mint    PublicKey `db:"mint"`
	Source  PublicKey `db:"source"`
	Wallet  PublicKey `db:"wallet"`
	*InstructionMetadata
}

func (s *AssociatedTokenAccountCreate) Table() string {
	return SPL_ASSOCIATED_TOKEN_ACCOUNT_CREATE
}

func (s *AssociatedTokenAccountCreate) Filter(k *keycache.Keycache) bool {
	return k.Contains(s.Mint.PublicKey())
}
