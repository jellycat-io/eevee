package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	TabSize int `toml:"tab_size"`
}

func GetConfig() Config {
	var config Config

	if _, err := toml.DecodeFile("eevee.toml", &config); err != nil {
		log.Fatal(err)
	}

	return config
}
