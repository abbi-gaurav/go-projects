package config

import (
	"github.com/vrischmann/envconfig"
	"log"
)

//AppConfig package variable to hold the Application config after it was loaded
var appConfig Config

//Config struct to hold the app config
type Config struct {
	CatalogFilePath string `envconfig:"default=/etc/sample-broker/catalog.yaml"`
	UserName        string
	Password        string
	Namespace       string `envconfig:"default=default"`
	ClusterDomain   string
}

//InitConfig initializes the AppConfig
func InitConfig() error {
	appConfig = Config{}
	err := envconfig.Init(&appConfig)
	log.Printf("Config %+v", appConfig)
	return err
}

//AppConfig returns the current AppConfig
func AppConfig() Config {
	return appConfig
}
