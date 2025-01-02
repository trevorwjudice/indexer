package spl

import (
	"fmt"
	"indexer/src/db/db_types"
	"indexer/src/util/solana/transactions"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
)

func PopulateInitializeAccount(r *transactions.Reader, flatIndex uint8, inst *token.InitializeAccount) (i *db_types.SplInitializeAccount, err error) {
	return &db_types.SplInitializeAccount{
		Owner:               db_types.PublicKey(inst.GetOwnerAccount().PublicKey),
		Mint:                db_types.PublicKey(inst.GetMintAccount().PublicKey),
		Account:             db_types.PublicKey(inst.GetAccount().PublicKey),
		InstructionMetadata: db_types.PopulateMetadata(r, flatIndex),
	}, nil
}

func PopulateInitializeAccount2(r *transactions.Reader, flatIndex uint8, inst *token.InitializeAccount2) (i *db_types.SplInitializeAccount, err error) {
	return &db_types.SplInitializeAccount{
		Owner:               db_types.PublicKey(*inst.Owner),
		Mint:                db_types.PublicKey(inst.GetMintAccount().PublicKey),
		Account:             db_types.PublicKey(inst.GetAccount().PublicKey),
		InstructionMetadata: db_types.PopulateMetadata(r, flatIndex),
	}, nil
}

func PopulateInitializeAccount3(r *transactions.Reader, flatIndex uint8, inst *token.InitializeAccount3) (i *db_types.SplInitializeAccount, err error) {
	return &db_types.SplInitializeAccount{
		Owner:               db_types.PublicKey(*inst.Owner),
		Mint:                db_types.PublicKey(inst.GetMintAccount().PublicKey),
		Account:             db_types.PublicKey(inst.GetAccount().PublicKey),
		InstructionMetadata: db_types.PopulateMetadata(r, flatIndex),
	}, nil
}

func PopulateTransfer(r *transactions.Reader, flatIndex uint8, inst *token.Transfer, prev []db_types.ParsedInstruction) (t *db_types.SplTransfer, err error) {
	t = &db_types.SplTransfer{
		Authority:           db_types.PublicKey(inst.GetOwnerAccount().PublicKey),
		Source:              db_types.PublicKey(inst.GetSourceAccount().PublicKey),
		Destination:         db_types.PublicKey(inst.GetDestinationAccount().PublicKey),
		Amount:              *inst.Amount,
		InstructionMetadata: db_types.PopulateMetadata(r, flatIndex),
	}
	t.Mint, err = db_types.ToPublicKeyErr(findMintAddress(r, prev, t.Source.PublicKey(), t.Destination.PublicKey()))
	if err != nil {
		return nil, err
	}
	return t, nil
}

func PopulateTransferChecked(r *transactions.Reader, flatIndex uint8, inst *token.TransferChecked, prev []db_types.ParsedInstruction) (t *db_types.SplTransfer, err error) {
	t = &db_types.SplTransfer{
		Authority:           db_types.PublicKey(inst.GetOwnerAccount().PublicKey),
		Source:              db_types.PublicKey(inst.GetSourceAccount().PublicKey),
		Destination:         db_types.PublicKey(inst.GetDestinationAccount().PublicKey),
		Amount:              *inst.Amount,
		InstructionMetadata: db_types.PopulateMetadata(r, flatIndex),
	}
	t.Mint, err = db_types.ToPublicKeyErr(findMintAddress(r, prev, t.Source.PublicKey(), t.Destination.PublicKey()))
	if err != nil {
		return nil, err
	}
	return t, nil
}

func PopulateBurn(r *transactions.Reader, flatIndex uint8, inst *token.Burn) (i *db_types.SplBurn, err error) {
	return &db_types.SplBurn{
		Mint:                db_types.PublicKey(inst.GetMintAccount().PublicKey),
		Account:             db_types.PublicKey(inst.GetSourceAccount().PublicKey),
		Owner:               db_types.PublicKey(inst.GetOwnerAccount().PublicKey),
		Amount:              *inst.Amount,
		InstructionMetadata: db_types.PopulateMetadata(r, flatIndex),
	}, nil
}

func PopulateBurnChecked(r *transactions.Reader, flatIndex uint8, inst *token.BurnChecked) (i *db_types.SplBurn, err error) {
	return &db_types.SplBurn{
		Mint:                db_types.PublicKey(inst.GetMintAccount().PublicKey),
		Account:             db_types.PublicKey(inst.GetSourceAccount().PublicKey),
		Owner:               db_types.PublicKey(inst.GetOwnerAccount().PublicKey),
		Amount:              *inst.Amount,
		InstructionMetadata: db_types.PopulateMetadata(r, flatIndex),
	}, nil
}

func findMintAddress(r *transactions.Reader, prev []db_types.ParsedInstruction, source, destination solana.PublicKey) (solana.PublicKey, error) {
	// First, check balance deltas to see if the source or destination account is referenced. This will work in all cases that transfers
	// from source and destination accounts do not result in 0 balances (almost all transfers).
	mint, err := r.FindTransferMintAddress(source, destination)
	if err == nil {
		return mint, nil
	}

	// If the above step does not succeed due to an edge case, look through previous instructions for a reference to the mint address.
	for _, inst := range prev {
		switch x := inst.(type) {
		case *db_types.SplInitializeAccount:
			if x.Account.PublicKey().Equals(source) || x.Account.PublicKey().Equals(destination) {
				return x.Mint.PublicKey(), nil
			}
		case *db_types.AssociatedTokenAccountCreate:
			if x.Account.PublicKey().Equals(source) || x.Account.PublicKey().Equals(destination) {
				return x.Mint.PublicKey(), nil
			}
		}
	}

	return *new(solana.PublicKey), fmt.Errorf("mint not found")
}
