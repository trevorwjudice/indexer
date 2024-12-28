package main

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
