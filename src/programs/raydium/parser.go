package raydium

import (
	"fmt"
	"indexer/src/indexer"
	"indexer/src/programs"
	"indexer/src/util/solana/transactions"

	"github.com/gagliardetto/solana-go"
)

var RAYDIUM_AMM_V4_PROGRAM_ID solana.PublicKey = solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")

type RaydiumInstructionParser struct{}

var _ indexer.InstructionParser = &RaydiumInstructionParser{}

func (r *RaydiumInstructionParser) ParseTransaction(reader *transactions.Reader) ([]indexer.ParsedInstruction, error) {
	if !reader.CheckAccountKeysContains(RAYDIUM_AMM_V4_PROGRAM_ID) {
		// Exit early if raydium v4 program id not in transaction account keys
		return nil, nil
	}

	flattened := reader.GetFlattenedInstructions()
	var res []indexer.ParsedInstruction
	for ind, inst := range flattened {
		programId, err := reader.GetAccountAtIndex(inst.ProgramIDIndex)
		if err != nil {
			return nil, err
		}

		if !programId.Equals(RAYDIUM_AMM_V4_PROGRAM_ID) {
			continue
		}

		// In this case, the instruction is a raydium instruction, so we should parse it.
		parsed, err := r.ParseInstruction(reader, uint8(ind))
		if err != nil {
			return nil, err
		}

		res = append(res, parsed)
	}

	return res, nil
}

func (r *RaydiumInstructionParser) ParseInstruction(reader *transactions.Reader, flatIndex uint8) (indexer.ParsedInstruction, error) {
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
		return nil, nil
		// return PopulateInitialize(reader, flatIndex)
	case InstructionInitialize2Discriminator:
		// return PopulateInitialize2(reader, flatIndex)
	case InstructionMonitorStepDiscriminator:
		// return PopulateMonitorStep(reader, flatIndex)
	case InstructionAddLiquidityDiscriminator:
		// return PopulateAddLiquidity(reader, flatIndex)
	case InstructionRemoveLiquidityDiscriminator:
		// return PopulateRemoveLiquidity(reader, flatIndex)
	case InstructionWithdrawPnlDiscriminator:
		// return PopulateWithdrawPnl(reader, flatIndex)
	case InstructionSwapExactAmountInDiscriminator:
		return PopulateSwapExactAmountIn(reader, flatIndex)
	case InstructionSwapExactAmountOutDiscriminator:
		return PopulateSwapExactAmountOut(reader, flatIndex)
	default:
		fmt.Println(reader.GetSignature())
		return nil, fmt.Errorf("unknown discriminator: %d", discriminator)
	}
	return nil, nil
}
