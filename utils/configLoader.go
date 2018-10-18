// Copyright 2018 OVO Technology
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
