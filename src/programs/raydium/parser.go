package raydium

import (
	"fmt"
	"indexer/src/db/db_types"
	"indexer/src/programs"
	"indexer/src/util/solana/transactions"
)

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
	case InstructionInitializeDiscriminator:
		panic("not implemented")
		// PopulateInitialize(reader, flatIndex)
	case InstructionInitialize2Discriminator:
		init2, err := PopulateInitialize2(reader, flatIndex)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{init2}, nil
	case InstructionMonitorStepDiscriminator:
		// Ignore
		return nil, nil
	case InstructionAddLiquidityDiscriminator:
		addLiquidity, err := PopulateAddLiquidity(reader, flatIndex)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{addLiquidity}, nil
	case InstructionRemoveLiquidityDiscriminator:
		removeLiqudiity, err := PopulateRemoveLiquidity(reader, flatIndex)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{removeLiqudiity}, nil
	case InstructionWithdrawPnlDiscriminator:
		// Ignore
		return nil, nil
	case InstructionSwapExactAmountInDiscriminator:
		swap, err := PopulateSwapExactAmountIn(reader, flatIndex)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{swap}, nil
	case InstructionSwapExactAmountOutDiscriminator:
		swap, err := PopulateSwapExactAmountOut(reader, flatIndex)
		if err != nil {
			return nil, err
		}
		return []db_types.ParsedInstruction{swap}, nil
	default:
		fmt.Println(reader.GetSignature())
		return nil, fmt.Errorf("unknown discriminator: %d", discriminator)
	}
}
