package utils

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	Order []string
}

const (
	prefFilename = "orderPref"
	envPrefix    = "helm_bulk"
)

//OrderPref returns a slice containing the preferred ordering of Release names.
// If it doesn't find any defined, it returns an empty slice.
func OrderPref(configDir string) (releaseOrderPref []string) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)
	viper.SetConfigName(prefFilename)
	viper.AddConfigPath(configDir)
	viper.ReadInConfig()
	var c config
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Println(err)
	}
	releaseOrderPref = c.Order
	return
}
