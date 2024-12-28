package pumpfun

import (
	"fmt"
	"indexer/src/indexer"
	"indexer/src/programs"
	"indexer/src/util/solana/transactions"
)

type PumpFunInstructionParser struct{}

var _ indexer.InstructionParser = &PumpFunInstructionParser{}

func (r *PumpFunInstructionParser) ParseTransaction(reader *transactions.Reader) ([]indexer.ParsedInstruction, error) {
	if !reader.CheckAccountKeysContains(PUMPFUN_PROGRAM_ID) {
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

		if !programId.Equals(PUMPFUN_PROGRAM_ID) {
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

func (r *PumpFunInstructionParser) ParseInstruction(reader *transactions.Reader, flatIndex uint8) (indexer.ParsedInstruction, error) {
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
		return PopulateCreate(reader, flatIndex)
	case InstructionBuyDiscriminator:
		return PopulateBuy(reader, flatIndex)
	case InstructionSellDiscriminator:
		return PopulateSell(reader, flatIndex)
	default:
		fmt.Println(reader.GetSignature())
		return nil, fmt.Errorf("unknown discriminator: %d", discriminator)
	}
}
