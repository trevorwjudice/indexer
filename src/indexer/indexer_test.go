package indexer_test

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigdotenv"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

type IndexerConfig struct {
	DB_HOST string `yaml:"DB_HOST"`
	DB_NAME string `yaml:"DB_NAME"`
	DB_PORT int    `yaml:"DB_PORT"`
	DB_USER string `yaml:"DB_USER"`
	DB_PW   string `yaml:"DB_PW"`
}

func GetConfig(configPath string) (*IndexerConfig, error) {
	c := &IndexerConfig{}

	loader := aconfig.LoaderFor(c, aconfig.Config{
		Files:              []string{configPath},
		AllowUnknownEnvs:   true,
		AllowUnknownFields: true,
		AllowUnknownFlags:  true,
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
			".env":  aconfigdotenv.New(),
		},
	})

	if err := loader.Load(); err != nil {
		return nil, err
	}

	return c, nil
}

// var cl = rpc.New("https://mainnet.helius-rpc.com/?api-key=c7c3e65e-8a33-4422-91e0-33b9d4764cca")
// var cl = rpc.New("https://mainnet.helius-rpc.com/?api-key=050a9592-0782-4c17-ab90-7e6aef937356")
// var ctx = context.TODO()
// var s db.Session
// var i *indexer.Indexer

func init() {
	// c, err := GetConfig("/Users/trevorjudice/Desktop/terminal/indexer/config/local.yaml")
	// if err != nil {
	// 	panic(err)
	// }

	// pgConnRW := postgresql.ConnectionURL{
	// 	Host:     c.DB_HOST,
	// 	Database: c.DB_NAME,
	// 	User:     c.DB_USER,
	// 	Password: c.DB_PW,
	// 	Options: map[string]string{
	// 		"sslmode": "prefer",
	// 	},
	// }

	// s, err = postgresql.Open(pgConnRW)
	// if err != nil {
	// 	panic(err)
	// }

	// i = indexer.NewIndexer(cl, s)

	// i.AddProgramParser(&raydium.RaydiumInstructionParser{})
}

// func TestParseTrans(t *testing.T) {
// validSlots, err := cl.GetBlocksWithLimit(ctx, 261717108, 1000, rpc.CommitmentConfirmed)
// for err != nil {
// 	validSlots, err = cl.GetBlocksWithLimit(ctx, 261717108, 1000, rpc.CommitmentConfirmed)
// 	if len(*validSlots) == 0 {
// 		err = errors.New("a")
// 	}
// }

// 	wg := errgroup.Group{}
// 	wg.SetLimit(200)
// 	for _, s := range *validSlots {
// 		s := s
// 		wg.Go(func() error {
// 			s := s
// 			_, err := i.GetSlot(ctx, s)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			return nil
// 		})
// 	}

// 	if err := wg.Wait(); err != nil {
// 		t.Fatal(err)
// 	}
// }
