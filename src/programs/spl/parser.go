package spl

import (
	"fmt"
	"indexer/src/db/db_types"
	"indexer/src/programs"
	"indexer/src/util/solana/transactions"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
)

var TOKEN_PROGRAM_ID solana.PublicKey = token.ProgramID

func ParseInstruction(r *transactions.Reader, flatIndex uint8, prev []db_types.ParsedInstruction) ([]db_types.ParsedInstruction, error) {
	inst, err := r.GetInstructionAtFlattenedIndex(flatIndex)
	if err != nil {
		return nil, err
	}
	disc, err := programs.GetInstructionDiscriminator(inst.Data)
	if err != nil {
		return nil, err
	}
	// Skip unrecognized instructions
	if disc == 21 || disc == 22 {
		return nil, nil
	}

	accounts, err := r.GetAccountMeta(inst)
	if err != nil {
		return nil, err
	}

	decoded, err := token.DecodeInstruction(accounts, inst.Data)
	if err != nil {
		fmt.Println("error parsing token instruction", err, r.GetSignature().String())
		return nil, nil
	}

	switch x := decoded.Impl.(type) {
	case *token.Transfer:
		transfer, err := PopulateTransfer(r, flatIndex, x, prev)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{transfer}, nil
	case *token.TransferChecked:
		transfer, err := PopulateTransferChecked(r, flatIndex, x, prev)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{transfer}, nil
	case *token.InitializeAccount:
		initializeAccount, err := PopulateInitializeAccount(r, flatIndex, x)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{initializeAccount}, nil
	case *token.InitializeAccount2:
		initializeAccount, err := PopulateInitializeAccount2(r, flatIndex, x)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{initializeAccount}, nil
	case *token.InitializeAccount3:
		initializeAccount, err := PopulateInitializeAccount3(r, flatIndex, x)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{initializeAccount}, nil
	case *token.Burn:
		burn, err := PopulateBurn(r, flatIndex, x)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{burn}, nil
	case *token.BurnChecked:
		burn, err := PopulateBurnChecked(r, flatIndex, x)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{burn}, nil
	}
	return nil, nil
}
