package pumpfun

import (
	"indexer/src/programs"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

var InstructionCreateDiscriminator programs.InstructionDiscriminator = 24
var InstructionBuyDiscriminator programs.InstructionDiscriminator = 102
var InstructionSellDiscriminator programs.InstructionDiscriminator = 51
var InstructionSetParamsDiscriminator programs.InstructionDiscriminator = 27

type InstructionCreate struct {
	Discriminator [8]byte
	Name          bin.SafeString
	Symbol        bin.SafeString
	URI           bin.SafeString
}

type InstructionBuy struct {
	Discriminator [8]byte
	Amount        uint64
	MaxSolCost    uint64
}

type InstructionSell struct {
	Discriminator [8]byte
	Amount        uint64
	MinSolOutput  uint64
}

type InstructionSetParams struct {
	Discriminator               [8]byte
	FeeRecipient                solana.PublicKey
	InitialVirtualTokenReserves uint64
	InitialVirtualSolReserves   uint64
	InitialRealTokenReserves    uint64
	TokenTotalSupply            uint64
	FeeBasisPoints              uint64
}
