package main

import (
	"electrumx-proxy-go/common/log"
	"electrumx-proxy-go/config"
	"electrumx-proxy-go/router"
	"electrumx-proxy-go/ws"
)

func main() {
	config.InitConf()

	log.InitLog(config.Conf.LogLevel,
		config.Conf.LogPath, log.Stdout)

	log.Infof("current version:%s", config.Version)

	if config.Conf.WssEnable {
		go ws.InitWebSocket(config.Conf.ElectrumxServer)
	}
	api := router.InitMasterRouter()
	err := api.Run(config.Conf.ServerAddress)
	if err != nil {
		log.Fatal(err)
		return
	}
}
