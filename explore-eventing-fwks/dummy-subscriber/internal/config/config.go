package config

import (
	"github.com/vrischmann/envconfig"
	"log"
)

type Options struct {
	LogRequest bool `envconfig:"default=false"`
}

var appConfig *Options

func InitConfig() {
	appConfig = &Options{}
	err := envconfig.Init(appConfig)
	if err != nil {
		log.Panicln(err)
	}
	log.Printf("config %+v\n", appConfig)
}

func AppConfig() *Options {
	return appConfig
}
