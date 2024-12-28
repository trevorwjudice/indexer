package solutils

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
)

func GetInstructionAccounts(accountIndicies []uint16, accounts []solana.PublicKey) ([]solana.PublicKey, error) {
	res := make([]solana.PublicKey, 0, len(accountIndicies))
	for _, ind := range accountIndicies {
		if int(ind) >= len(accounts) {
			return nil, fmt.Errorf("out of bounds")
		}
		res = append(res, accounts[ind])
	}
	return res, nil
}
