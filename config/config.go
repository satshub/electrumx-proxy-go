package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

var Version = "dev-dirty"

type Config struct {
	WssEnable       bool
	ElectrumxServer string
	ServerAddress   string
	LogLevel        int
	LogPath         string
}

var Conf Config

func InitConf() {
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		log.Fatal(err)
	}
}
