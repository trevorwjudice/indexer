package queries_test

import (
	"context"
	"fmt"
	"indexer/src/db/queries"
	"log"
	"os"
	"testing"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"gopkg.in/yaml.v3"
)

var s db.Session
var ctx context.Context = context.TODO()

type IndexerConfig struct {
	DB_HOST string `yaml:"DB_HOST"`
	DB_NAME string `yaml:"DB_NAME"`
	DB_PORT int    `yaml:"DB_PORT"`
	DB_USER string `yaml:"DB_USER"`
	DB_PW   string `yaml:"DB_PW"`
}

func GetConfig(configPath string) (*IndexerConfig, error) {
	c := &IndexerConfig{}

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c, nil
}

func init() {
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

	s, err = postgresql.Open(pgConnRW)
	if err != nil {
		panic(err)
	}

}

func TestGetPumpFunCreate(t *testing.T) {
	creates, err := queries.GetPumpFunCreateInstructions(ctx, s)
	if err != nil {
		t.Fatal(err)
	}

	for _, create := range creates {
		fmt.Println(create.Mint.PublicKey().String(), create.Name, create.Signature.Signature().String(), create.Slot)
		fmt.Printf("%+v\n", create)
	}
}

func TestGetPumpFunSwaps(t *testing.T) {
	swaps, err := queries.GetPumpFunSwaps(ctx, s)
	if err != nil {
		t.Fatal(err)
	}

	for _, swap := range swaps {
		fmt.Println(swap.Mint.PublicKey().String(), swap.Maker.PublicKey().String(), swap.Signature.Signature().String(), swap.Slot)
		fmt.Printf("%+v\n", swap)
	}
}

func TestGetRaydiumSwaps(t *testing.T) {
	swaps, err := queries.GetRaydiumSwaps(ctx, s, 10)
	if err != nil {
		t.Fatal(err)
	}

	for _, swap := range swaps {
		fmt.Println(swap.Signature.Signature().String(), swap.PoolIdentifier.PublicKey().String(), swap.Maker.PublicKey().String(), swap.Slot)
		fmt.Println(swap.AmountBase, swap.AmountQuote)
		fmt.Printf("%+v\n", swap)
	}
}
