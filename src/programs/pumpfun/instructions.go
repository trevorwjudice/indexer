package pumpfun

import (
	"indexer/src/programs"

	bin "github.com/gagliardetto/binary"
)

var InstructionCreateDiscriminator programs.InstructionDiscriminator = 24
var InstructionBuyDiscriminator programs.InstructionDiscriminator = 102
var InstructionSellDiscriminator programs.InstructionDiscriminator = 51

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
