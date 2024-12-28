package main

import (
	"context"
	"fmt"
	"indexer/src/indexer"
	"indexer/src/programs/pumpfun"
	"indexer/src/programs/raydium"
	"log"
	"os"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/upper/db/v4/adapter/postgresql"
)

var START_SLOT uint64 = 241768575

var ctx = context.TODO()

func main() {
	conf, err := GetConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		panic(err)
	}

	log.Println(os.Getenv("CONFIG_PATH"))

	pgConnRW := postgresql.ConnectionURL{
		Host:     conf.DB_HOST,
		Database: conf.DB_NAME,
		User:     conf.DB_USER,
		Password: conf.DB_PW,
		Options: map[string]string{
			"sslmode": "prefer",
		},
	}

	s, err := postgresql.Open(pgConnRW)
	if err != nil {
		panic(err)
	}

	evs, err := pumpfun.GetCreateInstructions(ctx, s)
	if err != nil {
		panic(err)
	}

	for _, ev := range evs {
		fmt.Println(ev.Name, ev.Signature, ev.Mint)
	}

	panic("a")

	_ = s

	// cl := rpc.New("https://winter-ultra-valley.solana-mainnet.quiknode.pro/a8e180434c2568f03f0159538c783447816e4ded")
	cl := rpc.New("https://mainnet.helius-rpc.com/?api-key=050a9592-0782-4c17-ab90-7e6aef937356")
	i := indexer.NewIndexer(cl, s)
	i.AddProgramParser(&raydium.RaydiumInstructionParser{})
	i.AddProgramParser(&pumpfun.PumpFunInstructionParser{})
	err = i.ScheduleRange(ctx, START_SLOT, 300000000)
	if err != nil {
		panic(err)
	}

	err = i.Run(ctx)
	if err != nil {
		panic(err)
	}
}
