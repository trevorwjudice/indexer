package instructions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gagliardetto/solana-go"
)

type StackTracePath struct {
	programId solana.PublicKey
	index     int
	path      []StackTracePath
}

func CreatePath(index int, programId solana.PublicKey) StackTracePath {
	return StackTracePath{
		index: index,
	}
}

func NewStackTracePath(path []StackTracePath) (StackTracePath, error) {
	if len(path) == 0 {
		return StackTracePath{}, fmt.Errorf("path length cannot be 0")
	}
	return StackTracePath{
		path: path,
	}, nil
}

func (s *StackTracePath) GetTopLevelProgramId() solana.PublicKey {
	return s.path[0].programId
}

func (s *StackTracePath) GetTopLevelInstructionIndex() int {
	return s.path[0].index
}

func (s *StackTracePath) GetProgramId() solana.PublicKey {
	return s.path[len(s.path)-1].programId
}

func (s *StackTracePath) GetStackDepth() int {
	return len(s.path)
}

func (s *StackTracePath) CopyPath() []StackTracePath {
	dst := make([]StackTracePath, len(s.path))
	copy(s.path, dst)
	return dst
}

func (s *StackTracePath) CreateNextSibling(programId solana.PublicKey) (StackTracePath, error) {
	sPath := s.CopyPath()
	sPath[len(sPath)-1].programId = programId
	sPath[len(sPath)-1].index += 1
	return NewStackTracePath(sPath)
}

func (s *StackTracePath) CreateChild(programId solana.PublicKey, index int) StackTracePath {
	c := StackTracePath{
		path:      s.CopyPath(),
		programId: programId,
		index:     index,
	}
	c.path = append(c.path, c)
	return c
}

func (s *StackTracePath) Equals(s2 StackTracePath) bool {
	if s.GetStackDepth() != s2.GetStackDepth() {
		return false
	}

	for i, sib := range s.path {
		sib2 := s2.path[i]
		if sib.index != sib2.index {
			return false
		}
		if !sib.programId.Equals(sib2.programId) {
			return false
		}
	}

	return true
}

func (s *StackTracePath) GetInstructionIdentifier() string {
	indices := make([]string, 0, len(s.path))
	for _, p := range s.path {
		indices = append(indices, strconv.FormatInt(int64(p.index), 10))
	}
	return "#" + strings.Join(indices, ".")
}

func (s *StackTracePath) String() string {
	if s.GetStackDepth() == 1 {
		return fmt.Sprintf("#%d (%s)", s.GetTopLevelInstructionIndex()+1, s.GetTopLevelProgramId())
	}

	stringified := make([]string, 0, len(s.path))
	for _, p := range s.path {
		stringified = append(stringified, fmt.Sprintf("#%d (%s)", p.index+1, p.programId))
	}
	return strings.Join(stringified, " -> ")
}
