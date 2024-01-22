package config

import (
	"log"

	"github.com/atanurdemir/gatekeeper/src/models"
	"github.com/spf13/viper"
)

var GatekeeperConfig models.AppConfig

func SetupConfig() {
	viper.SetConfigName("gates")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("src/config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	if err := viper.Unmarshal(&GatekeeperConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
}
