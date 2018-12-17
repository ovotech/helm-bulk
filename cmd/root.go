// Copyright © 2018 Chris Every <chris.every@ovoenergy.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
	loadCmd.Flags().BoolVarP(&disableTLS, "disable-tls", "t", false, "")
	rootCmd.PersistentFlags().StringVarP(&filePrefix, "fileprefix", "f",
		"helm-releases", "File prefix to use with a Load or Save command")
	rootCmd.PersistentFlags().StringVarP(&tlsKey, "tls-key-path", "", "", "")
	rootCmd.PersistentFlags().StringVarP(&tlsCert, "tls-cert-path", "", "", "")
	rootCmd.PersistentFlags().StringVarP(&caCert, "ca-cert-path", "", "", "")
	rootCmd.PersistentFlags().StringVarP(&tlsServerName, "tls-server-name", "", "", "")
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
