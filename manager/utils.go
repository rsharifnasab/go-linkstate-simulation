package main

import (
	"log"

	"github.com/spf13/viper"
)

func initLogger() {
	log.SetFlags(0)
	log.Println("Manager started")
}

func pnc(err error) {
	if err != nil {
		panic(err)
	}
}

func loadYaml(fileName string) *viper.Viper {
	config := viper.New()
	config.SetConfigName(fileName)
	config.AddConfigPath(".")
	pnc(config.ReadInConfig())
	return config
}
