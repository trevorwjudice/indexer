package db_types

import "indexer/src/indexer/keycache"

type RaydiumInitialize2 struct {
	PoolIdentifier       PublicKey `db:"pool"`
	Minter               PublicKey `db:"minter"`
	CoinMint             PublicKey `db:"coin_mint"`
	PoolCoinTokenAccount PublicKey `db:"pool_coin_token_account"`
	PcMint               PublicKey `db:"pc_mint"`
	PoolPcTokenAccount   PublicKey `db:"pool_pc_token_account"`
	LpMint               PublicKey `db:"lp_mint"`
	Nonce                uint8     `db:"nonce"`
	InitPcAmount         uint64    `db:"init_pc_amount"`
	InitCoinAmount       uint64    `db:"init_coin_amount"`
	LpAmount             uint64    `db:"lp_amount"`
	*InstructionMetadata
}

func (i *RaydiumInitialize2) Table() string {
	return RAYDIUM_V4_INITIALIZE2
}

func (c *RaydiumInitialize2) Filter(k *keycache.Keycache) bool {
	if k.Contains(c.CoinMint.PublicKey()) || k.Contains(c.PcMint.PublicKey()) {
		k.Add(c.LpMint.PublicKey())
		k.Add(c.PoolIdentifier.PublicKey())
		return true
	}
	return false
}

type RaydiumMonitorStep struct {
	PoolIdentifier   PublicKey `db:"pool"`
	PlanOrderLimit   uint16    `db:"plan_order_limit"`
	PlaceOrderLimit  uint16    `db:"place_order_limit"`
	CancelOrderLimit uint16    `db:"cancel_order_limit"`
	*InstructionMetadata
}

type RaydiumAddLiquidity struct {
	PoolIdentifier PublicKey `db:"pool"`
	Minter         PublicKey `db:"minter"`
	AmountBase     uint64    `db:"amount_base"`
	AmountQuote    uint64    `db:"amount_quote"`
	LpTokenAmount  uint64    `db:"lp_token_amount"`
	*InstructionMetadata
}

func (i *RaydiumAddLiquidity) Table() string {
	return RAYDIUM_V4_ADD_LIQUIDITY
}

func (i *RaydiumAddLiquidity) Filter(k *keycache.Keycache) bool {
	return k.Contains(i.PoolIdentifier.PublicKey())
}

type RaydiumRemoveLiquidity struct {
	PoolIdentifier PublicKey `db:"pool"`
	Owner          PublicKey `db:"owner"`
	AmountBase     uint64    `db:"amount_base"`
	AmountQuote    uint64    `db:"amount_quote"`
	LpTokenAmount  uint64    `db:"lp_token_amount"`
	*InstructionMetadata
}

func (i *RaydiumRemoveLiquidity) Table() string {
	return RAYDIUM_V4_REMOVE_LIQUIDITY
}

func (i *RaydiumRemoveLiquidity) Filter(k *keycache.Keycache) bool {
	return k.Contains(i.PoolIdentifier.PublicKey())
}

type RaydiumWithdrawPnl struct {
	PoolIdentifier PublicKey `db:"pool"`
	Owner          PublicKey `db:"owner"`
	AmountBase     uint64    `db:"amount_base"`
	AmountQuote    uint64    `db:"amount_quote"`
	*InstructionMetadata
}

type RaydiumSwap struct {
	PoolIdentifier PublicKey `db:"pool"`
	Maker          PublicKey `db:"maker"`
	AmountBase     int64     `db:"amount_base"`
	AmountQuote    int64     `db:"amount_quote"`
	*InstructionMetadata
}

func (s *RaydiumSwap) Table() string {
	return RAYDIUM_V4_SWAPS
}

func (s *RaydiumSwap) Filter(k *keycache.Keycache) bool {
	return k.Contains(s.PoolIdentifier.PublicKey())
}
