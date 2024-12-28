package programs

import "fmt"

type InstructionDiscriminator uint8

func GetInstructionDiscriminator(dat []byte) (InstructionDiscriminator, error) {
	if len(dat) == 0 {
		return InstructionDiscriminator(0), fmt.Errorf("data is 0 bytes long")
	}
	return InstructionDiscriminator(dat[0]), nil
}
