package transactions

type InstructionKey [40]byte

func (r *Reader) GetInstructionKey(flatIndex uint8) (k InstructionKey, err error) {
	inst, err := r.GetInstructionAtFlattenedIndex(flatIndex)
	if err != nil {
		return k, err
	}
	programId, err := r.GetAccountAtIndex(inst.ProgramIDIndex)
	if err != nil {
		return k, err
	}
	v := append(programId[:], inst.Data[0:8]...)
	copy(k[:], v)
	return k, nil
}
