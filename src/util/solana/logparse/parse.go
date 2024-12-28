package logparse

import (
	"fmt"
	"indexer/src/util/generic"
)

func ParseLogs(logs []string) ([]Log, error) {
	parsed := make([]Log, 0, len(logs))
	for _, l := range logs {
		res, err := parseLogLine(l)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, res)
	}
	return parsed, nil
}

type ParsedLogs struct {
	Logs      []Log
	Truncated bool
}

func ParseLogs2(logs []string) (*ParsedLogs, error) {
	flat := make([]Log, 0, len(logs))

	stack := generic.NewStack[Log]()

	for _, l := range logs {
		parsed, err := parseLogLine(l)
		if err != nil {
			return nil, err
		}

		switch t := parsed.(type) {
		case *Truncated:
			return &ParsedLogs{
				Logs:      flat,
				Truncated: true,
			}, nil
		case *ProgramInvoke:
			if t.Depth != int64(stack.Len())-1 {
				return nil, fmt.Errorf("invoke depth mismatch")
			}

		}
	}
	return nil, nil
}
