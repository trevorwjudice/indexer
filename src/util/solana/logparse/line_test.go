package logparse_test

import (
	"context"
	"indexer/src/indexer"
	"indexer/src/programs/raydium"
	"indexer/src/util/solana/logparse"
	"log"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
)

var cl = rpc.New("https://mainnet.helius-rpc.com/?api-key=050a9592-0782-4c17-ab90-7e6aef937356")
var ctx = context.TODO()

var i *indexer.Indexer

func init() {
	i = indexer.NewIndexer(cl, nil)

	i.AddProgramParser(&raydium.RaydiumInstructionParser{})
}

// func TestParseTrans(t *testing.T) {
// 	validSlots, err := cl.GetBlocksWithLimit(ctx, 241768281, 100, rpc.CommitmentConfirmed)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(*validSlots) == 0 {
// 		panic("err")
// 	}

// 	u := uint64(0)

// 	for _, s := range *validSlots {
// 		blk, err := cl.GetBlockWithOpts(ctx, s, &rpc.GetBlockOpts{
// 			MaxSupportedTransactionVersion: &u,
// 			Commitment:                     rpc.CommitmentConfirmed,
// 		})
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		for _, tx := range blk.Transactions {
// 			parsed, err := logparse.ParseLogs(tx.Meta.LogMessages)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			for _, p := range parsed {
// 				inv, ok := p.(*logparse.ProgramInvoke)
// 				if !ok {
// 					continue
// 				}

// 				fmt.Println(inv.ProgramId)
// 			}
// 		}
// 	}
// }

func TestFilterContext(t *testing.T) {
	validSlots, err := cl.GetBlocksWithLimit(ctx, 241768281, 100, rpc.CommitmentConfirmed)
	if err != nil {
		t.Fatal(err)
	}
	if len(*validSlots) == 0 {
		panic("err")
	}

	u := uint64(0)

	for _, s := range *validSlots {
		blk, err := cl.GetBlockWithOpts(ctx, s, &rpc.GetBlockOpts{
			MaxSupportedTransactionVersion: &u,
			Commitment:                     rpc.CommitmentConfirmed,
		})
		if err != nil {
			t.Fatal(err)
		}

		for _, tx := range blk.Transactions {
			if tx.Meta.Err != nil {
				continue
			}

			parsed, err := logparse.ParseLogs(tx.Meta.LogMessages)
			if err != nil {
				t.Fatal(err)
			}

			f := logparse.NewLogFilterer(parsed)

			typeFilter := "programLog"
			contextLogs, err := f.FilterProgramContext(raydium.RAYDIUM_AMM_V4_PROGRAM_ID, &typeFilter)
			if err != nil {
				t.Fatal(err)
			}
			if len(contextLogs) == 69 {
				log.Println(tx.MustGetTransaction().Signatures)
				panic("a")
			}
			for _, l := range contextLogs {
				log.Println(l.Raw())
			}
		}
	}
}
