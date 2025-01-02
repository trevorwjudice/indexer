package spl

import (
	"fmt"
	"indexer/src/db/db_types"
	"indexer/src/util/solana/transactions"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
)

var ASSOCIATED_TOKEN_ACCOUNT_PROGRAM_ID solana.PublicKey = associatedtokenaccount.ProgramID

func ParseAssociatedTokenAccountInstruction(r *transactions.Reader, flatIndex uint8, prev []db_types.ParsedInstruction) ([]db_types.ParsedInstruction, error) {
	inst, err := r.GetInstructionAtFlattenedIndex(flatIndex)
	if err != nil {
		return nil, err
	}

	accounts, err := r.GetAccountMeta(inst)
	if err != nil {
		return nil, err
	}

	decoded, err := associatedtokenaccount.DecodeInstruction(accounts, inst.Data)
	if err != nil {
		fmt.Println("error parsing token instruction", err, r.GetSignature().String())
		return nil, nil
	}

	createAccount, ok := decoded.Impl.(*associatedtokenaccount.Create)
	if !ok {
		return nil, nil
	}

	acc, _, err := solana.FindAssociatedTokenAddress(createAccount.Wallet, createAccount.Mint)
	if err != nil {
		return nil, err
	}

	acc2, err := r.GetAccountAtIndex(inst.Accounts[1])
	if err != nil {
		return nil, err
	}

	if acc2.Equals(acc) {
		return nil, fmt.Errorf("account mismatch: %s != %s in tx: %s", acc, acc2, r.GetSignature())
	}

	res := &db_types.AssociatedTokenAccountCreate{
		Account:             db_types.PublicKey(acc),
		Mint:                db_types.PublicKey(createAccount.Mint),
		Source:              db_types.PublicKey(createAccount.Payer),
		Wallet:              db_types.PublicKey(createAccount.Wallet),
		InstructionMetadata: db_types.PopulateMetadata(r, flatIndex),
	}

	return []db_types.ParsedInstruction{res}, nil
}
