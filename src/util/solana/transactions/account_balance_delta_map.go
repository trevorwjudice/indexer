package transactions

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	soltypes "github.com/gagliardetto/solana-go/rpc"
)

type TokenBalanceDelta struct {
	Mint          solana.PublicKey
	Owner         solana.PublicKey
	TokenAccount  solana.PublicKey
	BalanceBefore *soltypes.UiTokenAmount
	BalanceAfter  *soltypes.UiTokenAmount
}

type AccountBalanceDeltaMap struct {
	deltasByAccount map[solana.PublicKey]TokenBalanceDelta
}

func NewBalanceDeltaMap(preTokenBalances, postTokenBalances []soltypes.TokenBalance, accounts []solana.PublicKey) (*AccountBalanceDeltaMap, error) {
	res := &AccountBalanceDeltaMap{
		deltasByAccount: make(map[solana.PublicKey]TokenBalanceDelta),
	}

	for _, preBal := range preTokenBalances {
		delta := TokenBalanceDelta{
			Mint:          preBal.Mint,
			BalanceBefore: preBal.UiTokenAmount,
		}
		if int(preBal.AccountIndex) >= len(accounts) {
			return nil, fmt.Errorf("account index out of range")
		}
		delta.TokenAccount = accounts[preBal.AccountIndex]

		if preBal.Owner != nil {
			delta.Owner = *preBal.Owner
		}
		res.deltasByAccount[delta.TokenAccount] = delta
	}

	for _, postBal := range postTokenBalances {
		if int(postBal.AccountIndex) >= len(accounts) {
			return nil, fmt.Errorf("account index out of range")
		}
		tokenAccount := accounts[postBal.AccountIndex]
		delta, ok := res.deltasByAccount[tokenAccount]
		if !ok {
			delta = TokenBalanceDelta{
				Mint:         postBal.Mint,
				TokenAccount: tokenAccount,
			}
		}
		delta.BalanceAfter = postBal.UiTokenAmount
		if postBal.Owner != nil && delta.Owner == *new(solana.PublicKey) {
			delta.Owner = *postBal.Owner
		}

		res.deltasByAccount[delta.TokenAccount] = delta
	}

	for acc, delta := range res.deltasByAccount {
		if delta.BalanceAfter == nil {
			uiAmount := 0.0
			delta.BalanceAfter = &soltypes.UiTokenAmount{
				Amount:         "0",
				Decimals:       delta.BalanceBefore.Decimals,
				UiAmount:       &uiAmount,
				UiAmountString: "0.0",
			}
			res.deltasByAccount[acc] = delta
			continue
		}
		if delta.BalanceBefore == nil {
			uiAmount := 0.0
			delta.BalanceBefore = &soltypes.UiTokenAmount{
				Amount:         "0",
				Decimals:       delta.BalanceAfter.Decimals,
				UiAmount:       &uiAmount,
				UiAmountString: "0.0",
			}
			res.deltasByAccount[acc] = delta
		}
	}

	return res, nil
}

func (d *AccountBalanceDeltaMap) GetDelta(tokenAccount solana.PublicKey) (TokenBalanceDelta, bool) {
	delta, ok := d.deltasByAccount[tokenAccount]
	return delta, ok
}

func (d *AccountBalanceDeltaMap) GetMap() map[solana.PublicKey]TokenBalanceDelta {
	return d.deltasByAccount
}
