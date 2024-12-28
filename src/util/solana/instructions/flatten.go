package instructions

// func FlattenInstructionsWithStackTracePaths(tx *rpc.TransactionWithMeta) error {
// 	r, err := transactions.NewReader(tx)
// 	if err != nil {
// 		return err
// 	}

// 	topLevelInstructions := r.GetInstructions()

// 	parsedInstructions := make([]solana.CompiledInstruction, 0, len(topLevelInstructions))

// 	for topLevelInstrIndex, topLevelInstr := range topLevelInstructions {
// 		programId, err := r.GetAccountAtIndex(topLevelInstr.ProgramIDIndex)
// 		if err != nil {
// 			return err
// 		}
// 		parsedInstructions = append(parsedInstructions, topLevelInstr)

// 		trace, _ := NewStackTracePath([]StackTracePath{CreatePath(topLevelInstrIndex, programId)})

// 		innerInstructions := r.GetInnerInstructions(uint16(topLevelInstrIndex))
// 		_, _ = trace, innerInstructions
// 	}

// 	return nil
// }
