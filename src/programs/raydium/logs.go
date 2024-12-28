package raydium

// import (
// 	"encoding/base64"
// 	"fmt"
// 	"strings"

// 	bin "github.com/gagliardetto/binary"
// )

// type LogType uint8

// var (
// 	SwapBaseInLogType LogType = 3
// )

// func ParseLogLine(logLine string) (RaydiumLog, error) {
// 	words := strings.Split(logLine, " ")
// 	switch words[0] {
// 	case "Log":
// 		if words[1] == "truncated" {
// 			return nil, fmt.Errorf("truncated log")
// 		}
// 	case "Program":
// 	}
// 	return nil, nil
// }

// type RaydiumLog interface{}

// type SwapBaseInLog struct {
// 	AmountIn          uint64
// 	MinimumAmountOut  uint64
// 	Direction         uint64
// 	UserSourceAccount uint64
// 	PoolCoin          uint64
// 	PoolPc            uint64
// 	OutAmount         uint64
// }

// func ParseSwapBaseIn(logData []byte) (RaydiumLog, error) {
// 	r := &SwapBaseInLog{}
// 	return r, bin.NewBinDecoder(logData).Decode(r)
// }

// func Decode(data string) (RaydiumLog, error) {
// 	logDecoded, err := base64.StdEncoding.DecodeString(str)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(logDecoded) <= 1 {
// 		return nil, fmt.Errorf("invalid log length")
// 	}

// 	logType := LogType(logDecoded[0])
// 	logDecoded = logDecoded[1:]
// 	switch logType {
// 	case SwapBaseInLogType:
// 		return ParseSwapBaseIn(logDecoded)
// 	default:
// 		return nil, fmt.Errorf("invalid Log Type")
// 	}
// }
