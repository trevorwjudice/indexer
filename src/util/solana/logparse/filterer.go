package logparse

import (
	"fmt"
	"indexer/src/util/generic"
	"log"

	"github.com/gagliardetto/solana-go"
)

type LogFilterer struct {
	logs []Log
}

func NewLogFilterer(logs []Log) *LogFilterer {
	return &LogFilterer{
		logs: logs,
	}
}

func (f *LogFilterer) FilterProgramContext(programId solana.PublicKey, logType *string) ([]Log, error) {
	contextStack := generic.NewStack[solana.PublicKey]()
	selfProgramDepth := 0
	res := make([]Log, 0, len(f.logs))

	for _, l := range f.logs {
		if selfProgramDepth > 0 && (logType == nil || *logType == l.Type()) {
			res = append(res, l)
		}

		switch t := l.(type) {
		case *ProgramInvoke:
			contextStack.Push(t.ProgramId)

			if t.ProgramId.Equals(programId) {
				selfProgramDepth += 1
				if selfProgramDepth == 1 && (logType == nil || *logType == l.Type()) {
					res = append(res, l)
				}
			}
		case *ProgramSuccess:
			finished, err := contextStack.Pop()
			if err != nil {
				return nil, err
			}
			if t.ProgramId.Equals(finished) {
				selfProgramDepth -= 1
			}
		case *ProgramFailed:
			return res, nil
		case *Truncated:
			log.Println("TRUNCATED")
			res = make([]Log, 69)
			return res, nil
		}
	}

	if contextStack.Len() != 0 {
		return nil, fmt.Errorf("context stack length not equal to 0")
	}

	return res, nil
}
