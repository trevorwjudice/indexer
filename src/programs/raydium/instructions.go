package raydium

import (
	"indexer/src/programs"
)

var InstructionInitializeDiscriminator programs.InstructionDiscriminator = 0
var InstructionInitialize2Discriminator programs.InstructionDiscriminator = 1
var InstructionMonitorStepDiscriminator programs.InstructionDiscriminator = 2
var InstructionAddLiquidityDiscriminator programs.InstructionDiscriminator = 3
var InstructionRemoveLiquidityDiscriminator programs.InstructionDiscriminator = 4
var InstructionWithdrawPnlDiscriminator programs.InstructionDiscriminator = 7
var InstructionSwapExactAmountInDiscriminator programs.InstructionDiscriminator = 9
var InstructionSwapExactAmountOutDiscriminator programs.InstructionDiscriminator = 11

type InstructionInitialize struct {
	Discriminator programs.InstructionDiscriminator
	Nonce         uint8
	OpenTime      uint64
}

type InstructionInitialize2 struct {
	Discriminator  programs.InstructionDiscriminator
	Nonce          uint8
	OpenTime       uint64
	InitPcAmount   uint64
	InitCoinAmount uint64
}

type InstructionMonitorStep struct {
	Discriminator    programs.InstructionDiscriminator
	PlanOrderLimit   uint16
	PlaceOrderLimit  uint16
	CancelOrderLimit uint16
}

type InstructionPreInitialize struct {
	Discriminator programs.InstructionDiscriminator
	Nonce         uint8
}

type InstructionAddLiquidity struct {
	Discriminator programs.InstructionDiscriminator
	MaxCoinAmount uint64
	MaxPcAmount   uint64
	BaseSide      uint64
}

type InstructionRemoveLiquidity struct {
	Discriminator programs.InstructionDiscriminator
	Amount        uint64
}

type InstructionWithdrawPnl struct {
	Discriminator programs.InstructionDiscriminator
}

type InstructionSwapExactAmountIn struct {
	Discriminator uint8
	AmountIn      uint64
	MinAmountOut  uint64
}

type InstructionSwapExactAmountOut struct {
	Discriminator uint8
	MaxAmountIn   uint64
	AmountOut     uint64
}
