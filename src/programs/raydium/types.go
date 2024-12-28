package raydium

import (
	"indexer/src/db_types"
	"indexer/src/indexer"

	"github.com/gagliardetto/solana-go"
)

// type Initialize struct{}

type Initialize2 struct {
	PoolIdentifier       solana.PublicKey `db:"pool"`
	Minter               solana.PublicKey `db:"minter"`
	CoinMint             solana.PublicKey `db:"coin_mint"`
	PoolCoinTokenAccount solana.PublicKey `db:"pool_coin_token_account"`
	PcMint               solana.PublicKey `db:"pc_mint"`
	PoolPcTokenAccount   solana.PublicKey `db:"pool_pc_token_account"`
	LpMint               solana.PublicKey `db:"lp_mint"`
	Nonce                uint8            `db:"nonce"`
	InitPcAmount         uint64           `db:"init_pc_amount"`
	InitCoinAmount       uint64           `db:"init_coin_amount"`
	LpAmount             uint64           `db:"lp_amount"`
	*indexer.InstructionMetadata
}

func (i *Initialize2) Table() string {
	return "raydium_v4_initialize2"
}

type MonitorStep struct {
	PoolIdentifier   solana.PublicKey `db:"pool"`
	PlanOrderLimit   uint16           `db:"plan_order_limit"`
	PlaceOrderLimit  uint16           `db:"place_order_limit"`
	CancelOrderLimit uint16           `db:"cancel_order_limit"`
	*indexer.InstructionMetadata
}

type AddLiquidity struct {
	PoolIdentifier solana.PublicKey `db:"pool"`
	Minter         solana.PublicKey `db:"minter"`
	AmountBase     uint64           `db:"amount_base"`
	AmountQuote    uint64           `db:"amount_quote"`
	LpTokenAmount  uint64           `db:"lp_token_amount"`
	*indexer.InstructionMetadata
}

type RemoveLiquidity struct {
	PoolIdentifier solana.PublicKey `db:"pool"`
	Owner          solana.PublicKey `db:"owner"`
	AmountBase     uint64           `db:"amount_base"`
	AmountQuote    uint64           `db:"amount_quote"`
	LpTokenAmount  uint64           `db:"lp_token_amount"`
	*indexer.InstructionMetadata
}

type WithdrawPnl struct {
	PoolIdentifier solana.PublicKey `db:"pool"`
	Owner          solana.PublicKey `db:"owner"`
	AmountBase     uint64           `db:"amount_base"`
	AmountQuote    uint64           `db:"amount_quote"`
	*indexer.InstructionMetadata
}

type Swap struct {
	PoolIdentifier solana.PublicKey `db:"pool"`
	Maker          solana.PublicKey `db:"maker"`
	AmountBase     int64            `db:"amount_base"`
	AmountQuote    int64            `db:"amount_quote"`
	*indexer.InstructionMetadata
}

func (s *Swap) Table() string {
	return db_types.RAYDIUM_V4_SWAPS
}
