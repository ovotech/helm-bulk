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

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var local bool
var disableTLS bool
var filePrefix string
var tlsKey string
var tlsCert string
var caCert string
var tlsServerName string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "helm-bulk",
	Short: "Load or Save Releases from File to Cluster, or Cluster to File, respectively",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
	helmHome := os.Getenv("HELM_HOME")
	loadCmd.Flags().BoolVarP(&disableTLS, "disable-tls", "t", false, "")
	rootCmd.PersistentFlags().StringVarP(&filePrefix, "fileprefix", "f",
		"helm-releases", "File prefix to use with a Load or Save command")
	rootCmd.PersistentFlags().StringVarP(&tlsKey, "tls-key-path", "k",
		helmHome+"/key.pem", "Filepath of TLS key")
	rootCmd.PersistentFlags().StringVarP(&tlsCert, "tls-cert-path", "p",
		helmHome+"/cert.pem", "Filepath of TLS cert")
	rootCmd.PersistentFlags().StringVarP(&caCert, "ca-cert-path", "a",
		helmHome+"/ca.pem", "Filepath of CA cert")
	rootCmd.PersistentFlags().StringVarP(&tlsServerName, "tls-server-name", "s",
		"", "TLS server name")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".helm-bulk" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".helm-bulk")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// textFilename returns the text filename
func textFilename() (filename string) {
	filename = filePrefix + ".txt"
	return
}

// archiveFilename returns the archive filename
func archiveFilename() (filename string) {
	filename = filePrefix + ".tar.gz"
	return
}
