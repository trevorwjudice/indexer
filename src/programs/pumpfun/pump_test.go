package pumpfun_test

import (
	"context"
	"errors"
	"indexer/src/indexer"
	"indexer/src/programs/pumpfun"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/upper/db/v4"
)

var cl = rpc.New("https://mainnet.helius-rpc.com/?api-key=050a9592-0782-4c17-ab90-7e6aef937356")
var ctx = context.TODO()
var s db.Session
var i *indexer.Indexer

func init() {
	i = indexer.NewIndexer(cl, s)

	i.AddProgramParser(&pumpfun.PumpFunInstructionParser{})
}

// First slot with create: 241768575
func TestParseTrans(t *testing.T) {
	validSlots, err := cl.GetBlocksWithLimit(ctx, 261768575, 1000, rpc.CommitmentConfirmed)
	for err != nil {
		validSlots, err = cl.GetBlocksWithLimit(ctx, 261768575, 1000, rpc.CommitmentConfirmed)
		if len(*validSlots) == 0 {
			err = errors.New("a")
		}
	}

	for _, s := range *validSlots {
		_, err := i.GetSlot(ctx, s)
		if err != nil {
			t.Fatal(err)
		}
	}
}
