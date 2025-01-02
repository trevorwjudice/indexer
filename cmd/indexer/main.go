package main

import (
	"context"
	"indexer/src/db/db_types"
	"indexer/src/indexer"
	"indexer/src/programs/pumpfun"
	"indexer/src/programs/raydium"
	"indexer/src/programs/spl"
	"log"
	"os"

	"net/http"
	_ "net/http/pprof"

	"github.com/gagliardetto/solana-go"
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

	cl := rpc.New("https://mainnet.helius-rpc.com/?api-key=050a9592-0782-4c17-ab90-7e6aef937356")
	i := indexer.NewIndexer(cl, s, db_types.NewFetchWhitelistedAddressesFunc([]solana.PublicKey{raydium.SOL_USDC_POOL}))
	i.AddParser(pumpfun.PUMPFUN_PROGRAM_ID, pumpfun.ParseInstruction)
	i.AddParser(raydium.RAYDIUM_AMM_V4_PROGRAM_ID, raydium.ParseInstruction)
	i.AddParser(spl.ASSOCIATED_TOKEN_ACCOUNT_PROGRAM_ID, spl.ParseAssociatedTokenAccountInstruction)
	i.AddParser(spl.TOKEN_PROGRAM_ID, spl.ParseInstruction)
	err = i.ScheduleRange(ctx, START_SLOT, 300000000)
	if err != nil {
		panic(err)
	}

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	err = i.Run(ctx)
	if err != nil {
		panic(err)
	}
}
