package pumpfun

import (
	"fmt"
	"indexer/src/db/db_types"
	"indexer/src/programs"
	"indexer/src/util/solana/transactions"

	"github.com/gagliardetto/solana-go"
)

var PUMPFUN_PROGRAM_ID solana.PublicKey = solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")

func ParseInstruction(reader *transactions.Reader, flatIndex uint8, _ []db_types.ParsedInstruction) ([]db_types.ParsedInstruction, error) {
	inst, err := reader.GetInstructionAtFlattenedIndex(flatIndex)
	if err != nil {
		return nil, err
	}

	discriminator, err := programs.GetInstructionDiscriminator(inst.Data)
	if err != nil {
		return nil, err
	}

	switch discriminator {
	case InstructionCreateDiscriminator:
		create, err := PopulateCreate(reader, flatIndex)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{create}, nil
	case InstructionBuyDiscriminator:
		buy, err := PopulateBuy(reader, flatIndex)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{buy}, nil
	case InstructionSellDiscriminator:
		sell, err := PopulateSell(reader, flatIndex)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{sell}, nil
	case InstructionSetParamsDiscriminator:
		setParams, err := PopulateSetParams(reader, flatIndex)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{setParams}, nil
	default:
		fmt.Println(reader.GetSignature())
		return nil, fmt.Errorf("unknown discriminator: %d", discriminator)
	}
}
