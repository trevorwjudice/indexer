package indexer_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var cl = rpc.New("https://mainnet.helius-rpc.com/?api-key=050a9592-0782-4c17-ab90-7e6aef937356")
var ctx = context.TODO()

var trials = 30

func TestGetBlockEncodingJson(t *testing.T) {
	fmt.Println("================ JSON ================")

	validSlots := fetchValidSlots()
	timeTotal := 0.0
	for _, s := range validSlots {
		ts := time.Now()
		_, err := cl.GetBlockWithOpts(ctx, s, &rpc.GetBlockOpts{
			Commitment:                     rpc.CommitmentFinalized,
			MaxSupportedTransactionVersion: rpc.NewTransactionVersion(0),
		})
		if err != nil {
			t.Fatal(err)
		}

		elapsed := time.Since(ts).Seconds()
		timeTotal += elapsed

		fmt.Println("Time Taken:", elapsed)
	}

	fmt.Println("Total Time Taken: ", timeTotal)
	fmt.Println("Blocks Fetched: ", len(validSlots))
	fmt.Println("Average Fetch Time: ", timeTotal/float64(len(validSlots)))
}

func TestGetBlockEncodingBase58(t *testing.T) {
	fmt.Println("================ Base 58 ================")
	validSlots := fetchValidSlots()
	timeTotal := 0.0
	for _, s := range validSlots {
		ts := time.Now()
		_, err := cl.GetBlockWithOpts(ctx, s, &rpc.GetBlockOpts{
			Commitment:                     rpc.CommitmentFinalized,
			MaxSupportedTransactionVersion: rpc.NewTransactionVersion(0),
			Encoding:                       solana.EncodingBase58,
		})
		if err != nil {
			t.Fatal(err)
		}

		elapsed := time.Since(ts).Seconds()
		timeTotal += elapsed

		fmt.Println("Time Taken:", elapsed)
	}

	fmt.Println("Total Time Taken: ", timeTotal)
	fmt.Println("Blocks Fetched: ", len(validSlots))
	fmt.Println("Average Fetch Time: ", timeTotal/float64(len(validSlots)))
}

func TestGetBlockEncodingBase64(t *testing.T) {
	fmt.Println("================ Base 64 ================")
	validSlots := fetchValidSlots()
	timeTotal := 0.0
	for _, s := range validSlots {
		ts := time.Now()
		_, err := cl.GetBlockWithOpts(ctx, s, &rpc.GetBlockOpts{
			Commitment:                     rpc.CommitmentFinalized,
			MaxSupportedTransactionVersion: rpc.NewTransactionVersion(0),
			Encoding:                       solana.EncodingBase64,
		})
		if err != nil {
			t.Fatal(err)
		}

		elapsed := time.Since(ts).Seconds()
		timeTotal += elapsed

		fmt.Println("Time Taken:", elapsed)
	}

	fmt.Println("Total Time Taken: ", timeTotal)
	fmt.Println("Blocks Fetched: ", len(validSlots))
	fmt.Println("Average Fetch Time: ", timeTotal/float64(len(validSlots)))
}

func TestGetBlockEncodingBase64ZStd(t *testing.T) {
	fmt.Println("================ Base 64ZStd ================")
	validSlots := fetchValidSlots()
	timeTotal := 0.0
	for _, s := range validSlots {
		ts := time.Now()
		_, err := cl.GetBlockWithOpts(ctx, s, &rpc.GetBlockOpts{
			Commitment:                     rpc.CommitmentFinalized,
			MaxSupportedTransactionVersion: rpc.NewTransactionVersion(0),
			Encoding:                       solana.EncodingBase64Zstd,
		})
		if err != nil {
			t.Fatal(err)
		}

		elapsed := time.Since(ts).Seconds()
		timeTotal += elapsed

		fmt.Println("Time Taken:", elapsed)
	}

	fmt.Println("Total Time Taken: ", timeTotal)
	fmt.Println("Blocks Fetched: ", len(validSlots))
	fmt.Println("Average Fetch Time: ", timeTotal/float64(len(validSlots)))
}

func fetchValidSlots() []uint64 {
	validSlots, err := cl.GetBlocksWithLimit(ctx, 261717108, uint64(trials), rpc.CommitmentConfirmed)
	for err != nil {
		validSlots, err = cl.GetBlocksWithLimit(ctx, 261717108, uint64(trials), rpc.CommitmentConfirmed)
		if len(*validSlots) == 0 {
			err = errors.New("a")
		}
	}
	return *validSlots
}
