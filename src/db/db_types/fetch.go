package db_types

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/upper/db/v4"
)

var FETCH_ADDRESS_QUERY = "(SELECT mint AS whitelisted_address FROM %s) UNION (SELECT pool AS whitelisted_address FROM %s) UNION (SELECT lp_mint AS whitelisted_address FROM %s)"

func FetchWhitelistedAddresses(ctx context.Context, s db.Session) ([]solana.PublicKey, error) {
	type whitelisted_address struct {
		WhitelistedAddress PublicKey `db:"whitelisted_address"`
	}
	queryRes := []whitelisted_address{}
	err := s.SQL().IteratorContext(ctx, fmt.Sprintf("(SELECT mint AS whitelisted_address FROM %s) UNION (SELECT pool AS whitelisted_address FROM %s) UNION (SELECT lp_mint AS whitelisted_address FROM %s)", PUMP_FUN_CREATE, RAYDIUM_V4_INITIALIZE2, RAYDIUM_V4_INITIALIZE2)).All(&queryRes)
	if err != nil {
		return nil, err
	}

	res := make([]solana.PublicKey, 0, len(queryRes))
	for _, q := range queryRes {
		res = append(res, q.WhitelistedAddress.PublicKey())
	}
	return res, nil
}

func NewFetchWhitelistedAddressesFunc(constAddrs []solana.PublicKey) func(ctx context.Context, s db.Session) ([]solana.PublicKey, error) {
	return func(ctx context.Context, s db.Session) ([]solana.PublicKey, error) {
		type whitelisted_address struct {
			WhitelistedAddress PublicKey `db:"whitelisted_address"`
		}
		queryRes := []whitelisted_address{}
		err := s.SQL().IteratorContext(ctx, fmt.Sprintf("(SELECT mint AS whitelisted_address FROM %s) UNION (SELECT pool AS whitelisted_address FROM %s) UNION (SELECT lp_mint AS whitelisted_address FROM %s)", PUMP_FUN_CREATE, RAYDIUM_V4_INITIALIZE2, RAYDIUM_V4_INITIALIZE2)).All(&queryRes)
		if err != nil {
			return nil, err
		}

		res := make([]solana.PublicKey, 0, len(queryRes)+len(constAddrs))
		for _, q := range queryRes {
			res = append(res, q.WhitelistedAddress.PublicKey())
		}
		res = append(res, constAddrs...)
		return res, nil
	}
}
