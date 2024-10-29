package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	WssEnable       bool
	ElectrumxServer string
	ServerAddress   string
}

var Conf Config

func InitConf() {
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		log.Fatal(err)
	}
}
