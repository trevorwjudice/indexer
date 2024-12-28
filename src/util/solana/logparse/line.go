package logparse

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gagliardetto/solana-go"
)

var (
	ErrorEmptyLog               error = errors.New("logparse: empty log")
	ErrorInvalidProgramId       error = errors.New("logparse: invalid program id")
	ErrorIncorrectLineWordCount error = errors.New("logparse: incorrect line word count")
	ErrorUnexpectedLogFormat    error = errors.New("logparse: unexpected log format")
)

type Log interface {
	Type() string
	Raw() string
}

type Truncated struct {
	raw string
}

var _ Log = &Truncated{}

func (t *Truncated) Type() string {
	return "logTruncated"
}

func (t *Truncated) Raw() string {
	return t.raw
}

func ParseTruncateLog(raw string, words []string) (*Truncated, error) {
	if len(words) != 2 {
		return nil, ErrorIncorrectLineWordCount
	}

	if words[0] != "Log" || words[1] != "truncated" {
		return nil, ErrorUnexpectedLogFormat
	}

	return &Truncated{
		raw: raw,
	}, nil
}

type ProgramConsumption struct {
	raw string
}

var _ Log = &ProgramConsumption{}

func (t *ProgramConsumption) Type() string {
	return "programConsumption"
}

func (t *ProgramConsumption) Raw() string {
	return t.raw
}

func ParseProgramConsumption(raw string, words []string) (*ProgramConsumption, error) {
	if words[0] != "Program" || words[1] != "consumption:" {
		return nil, ErrorUnexpectedLogFormat
	}

	return &ProgramConsumption{
		raw: raw,
	}, nil
}

type ProgramInvoke struct {
	ProgramId solana.PublicKey
	Depth     int64
	raw       string
}

var _ Log = &ProgramInvoke{}

func (t *ProgramInvoke) Type() string {
	return "programInvoke"
}

func (t *ProgramInvoke) Raw() string {
	return t.raw
}

func ParseProgramInvoke(raw string, words []string) (*ProgramInvoke, error) {
	if words[0] != "Program" || words[2] != "invoke" {
		return nil, ErrorUnexpectedLogFormat
	}

	programId, err := solana.PublicKeyFromBase58(words[1])
	if err != nil {
		return nil, errors.New("logparse: failed to parse program id")
	}

	depth, err := parseDepth(words[3])
	if err != nil {
		return nil, errors.New("logparse: failed to parse depth")
	}

	return &ProgramInvoke{
		ProgramId: programId,
		Depth:     depth,
		raw:       raw,
	}, nil
}

type ProgramSuccess struct {
	ProgramId solana.PublicKey
	raw       string
}

var _ Log = &ProgramSuccess{}

func (t *ProgramSuccess) Type() string {
	return "programSuccess"
}

func (t *ProgramSuccess) Raw() string {
	return t.raw
}

func ParseProgramSuccess(raw string, words []string) (*ProgramSuccess, error) {
	if words[0] != "Program" || words[2] != "success" {
		return nil, ErrorUnexpectedLogFormat
	}
	programId, err := solana.PublicKeyFromBase58(words[1])
	if err != nil {
		return nil, fmt.Errorf("logparse: failed to parse program id")
	}

	return &ProgramSuccess{
		ProgramId: programId,
		raw:       raw,
	}, nil
}

type ProgramFailed struct {
	ProgramId solana.PublicKey
	Error     string
	raw       string
}

var _ Log = &ProgramFailed{}

func (t *ProgramFailed) Type() string {
	return "programFailed"
}

func (t *ProgramFailed) Raw() string {
	return t.raw
}

func ParseProgramFailed(raw string, words []string) (*ProgramFailed, error) {
	if words[0] != "Program" || words[2] != "failed:" {
		return nil, ErrorUnexpectedLogFormat
	}
	programId, err := solana.PublicKeyFromBase58(words[1])
	if err != nil {
		return nil, fmt.Errorf("logparse: failed to parse program id")
	}

	return &ProgramFailed{
		ProgramId: programId,
		Error:     strings.Join(words[3:], " "),
		raw:       raw,
	}, nil
}

type ProgramCompleteFailed struct {
	Error string
	raw   string
}

var _ Log = &ProgramCompleteFailed{}

func (t *ProgramCompleteFailed) Type() string {
	return "programcompleteFailed"
}

func (t *ProgramCompleteFailed) Raw() string {
	return t.raw
}

func ParseProgramCompleteFailed(raw string, words []string) (*ProgramCompleteFailed, error) {
	if !(words[1] == "failed" || words[2] == "to" || words[3] == "complete:") {
		return nil, ErrorUnexpectedLogFormat
	}
	return &ProgramCompleteFailed{
		Error: strings.Join(words[4:], " "),
		raw:   raw,
	}, nil
}

type ProgramLog struct {
	Message string
	raw     string
}

var _ Log = &ProgramLog{}

func (t *ProgramLog) Type() string {
	return "programLog"
}

func (t *ProgramLog) Raw() string {
	return t.raw
}

func ParseProgramLog(raw string, words []string) (*ProgramLog, error) {
	if words[0] != "Program" || words[1] != "log:" {
		return nil, ErrorUnexpectedLogFormat
	}
	return &ProgramLog{
		Message: strings.Join(words[2:], " "),
		raw:     raw,
	}, nil
}

type ProgramData struct {
	Data string
	raw  string
}

var _ Log = &ProgramData{}

func (t *ProgramData) Type() string {
	return "programData"
}

func (t *ProgramData) Raw() string {
	return t.raw
}

func ParseProgramData(raw string, words []string) (*ProgramData, error) {
	if words[0] != "Program" || words[1] != "data:" {
		return nil, ErrorUnexpectedLogFormat
	}
	return &ProgramData{
		Data: strings.Join(words[2:], " "),
		raw:  raw,
	}, nil
}

type ProgramConsumed struct {
	ProgramId solana.PublicKey
	Consumed  uint64
	Total     uint64
	raw       string
}

var _ Log = &ProgramConsumed{}

func (t *ProgramConsumed) Type() string {
	return "programConsumed"
}

func (t *ProgramConsumed) Raw() string {
	return t.raw
}

func ParseProgramConsumed(raw string, words []string) (*ProgramConsumed, error) {
	if words[0] != "Program" || words[2] != "consumed" {
		return nil, ErrorUnexpectedLogFormat
	}
	programId, err := solana.PublicKeyFromBase58(words[1])
	if err != nil {
		return nil, err
	}
	consumed, err := strconv.ParseUint(words[3], 10, 64)
	if err != nil {
		return nil, err
	}
	total, err := strconv.ParseUint(words[5], 10, 64)
	if err != nil {
		return nil, err
	}
	return &ProgramConsumed{
		ProgramId: programId,
		Consumed:  consumed,
		Total:     total,
		raw:       raw,
	}, nil
}

type ProgramReturn struct {
	ProgramId solana.PublicKey
	Result    string
	raw       string
}

var _ Log = &ProgramReturn{}

func (t *ProgramReturn) Type() string {
	return "programReturn"
}

func (t *ProgramReturn) Raw() string {
	return t.raw
}

func ParseProgramReturn(raw string, words []string) (*ProgramReturn, error) {
	if words[0] != "Program" || words[1] != "return:" {
		return nil, ErrorUnexpectedLogFormat
	}
	programId, err := solana.PublicKeyFromBase58(words[2])
	if err != nil {
		return nil, err
	}
	return &ProgramReturn{
		ProgramId: programId,
		Result:    strings.Join(words[3:], " "),
		raw:       raw,
	}, nil
}

type TransferInsufficientLamports struct {
	raw string
}

var _ Log = &TransferInsufficientLamports{}

func (t *TransferInsufficientLamports) Type() string {
	return "transferInsufficientLamports"
}

func (t *TransferInsufficientLamports) Raw() string {
	return t.raw
}

func ParseTransferInsufficientLamports(raw string, words []string) (*TransferInsufficientLamports, error) {
	if words[0] != "Transfer:" || words[1] != "insufficient" || words[2] != "lamports" {
		return nil, ErrorUnexpectedLogFormat
	}
	return &TransferInsufficientLamports{
		raw: raw,
	}, nil
}

type StakeProgramLog struct {
	raw string
}

var _ Log = &StakeProgramLog{}

func (t *StakeProgramLog) Type() string {
	return "stakeProgramLog"
}

func (t *StakeProgramLog) Raw() string {
	return t.raw
}

func ParseStakeProgramLog(raw string, words []string) (*StakeProgramLog, error) {
	if !(words[1] == "if" && (words[2] == "destination" || words[2] == "source") && words[3] == "stake" && words[4] == "is" && words[5] == "mergeable") {
		return nil, ErrorUnexpectedLogFormat
	}
	return &StakeProgramLog{
		raw: raw,
	}, nil
}

type Unknown struct {
	raw string
}

var _ Log = &StakeProgramLog{}

func (t *Unknown) Type() string {
	return "unknown"
}

func (t *Unknown) Raw() string {
	return t.raw
}

func ParseUnknown(raw string) (*Unknown, error) {
	return &Unknown{
		raw: raw,
	}, nil
}

func parseLogLine(line string) (Log, error) {
	if len(line) == 0 {
		return nil, fmt.Errorf("logparse: line length cannot be 0")
	}

	words := strings.Split(line, " ")

	switch words[0] {
	case "Log":
		if words[1] == "truncated" {
			return ParseTruncateLog(line, words)
		}
	case "Program":
		if words[2] == "invoke" {
			return ParseProgramInvoke(line, words)
		}
		if words[2] == "success" {
			return ParseProgramSuccess(line, words)
		}
		if words[2] == "failed:" {
			return ParseProgramFailed(line, words)
		}
		if words[1] == "failed" && words[2] == "to" && words[3] == "complete:" {
			return ParseProgramCompleteFailed(line, words)
		}
		if words[1] == "log:" {
			return ParseProgramLog(line, words)
		}
		if words[1] == "data:" {
			return ParseProgramData(line, words)
		}
		if words[1] == "return:" {

		}
		if words[2] == "consumed" {
			return ParseProgramConsumed(line, words)
		}
		if words[1] == "return:" {
			return ParseProgramReturn(line, words)
		}
	case "Transfer:":
		if words[1] == "insufficient" && words[2] == "lamports" {
			return ParseTransferInsufficientLamports(line, words)
		}
	case "Checking":
		// (words[1] === "if" && (words[2] === "destination" || words[2] === "source") && words[3] === "stake" && words[4] === "is" && words[5] === "mergeable")
		if words[1] == "if" && (words[2] == "destination" || words[2] == "source") && words[3] == "stake" && words[4] == "is" && words[5] == "mergeable" {
			return ParseStakeProgramLog(line, words)
		}
	case "Merging":
		if words[1] == "stake" && words[2] == "accounts" {
			return &StakeProgramLog{
				raw: line,
			}, nil
		}
	}
	return ParseUnknown(line)
}

func parseDepth(word string) (int64, error) {
	intStr := word[1 : len(word)-1]

	return strconv.ParseInt(intStr, 10, 64)
}
