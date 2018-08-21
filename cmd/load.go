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
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/release"
	//"k8s.io/helm/_proto/hapi/rudder"
)

var b64 = base64.StdEncoding

var magicGzip = []byte{0x1f, 0x8b, 0x08}

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("load called")
		load()
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func load() {
	dat, err := ioutil.ReadFile("helm-releases.txt")
	check(err)
	for _, splitString := range strings.Split(string(dat), ",") {
		rel, err := decodeRelease(splitString)
		check(err)
		fmt.Println(rel.Name)
		client := helm.NewClient(helm.Host("127.0.0.1:44134"))
		var disableHooks = true
		var releaseName = rel.GetName()
		var namespace = rel.GetNamespace()
		var reuseName = true
		var dryRun = true
		//var chartName = rel.GetChart().GetMetadata().GetName()
		//var overrides = []byte("key1=value1,key2=value2")
		var overrides = []byte(rel.GetChart().GetValues().String())

		// Options used in InstallRelease
		ops := []helm.InstallOption{
			helm.ValueOverrides(overrides),
			helm.InstallDryRun(dryRun),
			helm.ReleaseName(releaseName),
			helm.InstallReuseName(reuseName),
			helm.InstallDisableHooks(disableHooks),
		}
		client.InstallReleaseFromChart(rel.Chart, namespace, ops...)
	}
}

// decodeRelease decodes the bytes in data into a release
// type. Data must contain a base64 encoded string of a
// valid protobuf encoding of a release, otherwise
// an error is returned.
func decodeRelease(data string) (*release.Release, error) {
	// base64 decode string
	b, err := b64.DecodeString(data)
	if err != nil {
		return nil, err
	}

	// For backwards compatibility with releases that were stored before
	// compression was introduced we skip decompression if the
	// gzip magic header is not found
	if bytes.Equal(b[0:3], magicGzip) {
		r, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		b2, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		b = b2
	}

	var rls release.Release
	// unmarshal protobuf bytes
	if err := proto.Unmarshal(b, &rls); err != nil {
		return nil, err
	}
	return &rls, nil
}

//check panics if error is not nil
func check(e error) {
	if e != nil {
		panic(e.Error())
	}
}
