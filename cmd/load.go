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
	"io/ioutil"
	"log"
	"strings"

	"github.com/ovotech/helm-bulk/utils"
	"github.com/spf13/cobra"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/release"
)

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
		log.Println("helm-bulk load called")
		releases := releases()
		logReleasesFound(releases)
		load(releases)
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

//logReleasesFound log.Println's the names of releases, or a message indicating
//none were found
func logReleasesFound(releases []*release.Release) {
	if len(releases) > 0 {
		var buffer bytes.Buffer
		buffer.WriteString("Found Helm Releases to load:")
		buffer.WriteString("\n")
		for _, release := range releases {
			buffer.WriteString(release.GetName())
			buffer.WriteString("\n")
		}
		log.Println(buffer.String())
	} else {
		log.Println("No Helm Releases found to load")
	}

}

//releases decodes the release file and returns a slice of releases
func releases() (releases []*release.Release) {
	dat, err := ioutil.ReadFile("helm-releases.txt")
	panicCheck(err)
	for _, splitString := range strings.Split(string(dat), ",") {
		release, err := utils.DecodeRelease(splitString)
		panicCheck(err)
		releases = append(releases, release)
	}
	return
}

//load attempts to install each release in the provided slice.
//If an error is encountered in doing so, it logs the failure and skips to the
//next element in the slice
func load(releases []*release.Release) {
	for _, release := range releases {
		releaseName := release.GetName()
		log.Println("loading release:", releaseName)
		client := helm.NewClient(helm.Host("127.0.0.1:44134"))
		options := options(release)
		resp, err := client.InstallReleaseFromChart(release.Chart,
			release.GetNamespace(), options...)
		if err != nil {
			log.Println("loading of release:", releaseName, " failed:", err.Error())
			continue
		}
		log.Println(releaseName, "'helm install' response status:",
			resp.GetRelease().GetInfo().GetStatus().GetCode())
	}
}

//options creates and returns a slice of InstallOptions, of which some fields
//are grabbed from the provided release
func options(release *release.Release) (options []helm.InstallOption) {
	var disableHooks = true
	var releaseName = release.GetName()
	var reuseName = false
	var dryRun = true
	cv := release.GetChart().GetValues()
	var overrides = []byte(cv.Raw)
	options = []helm.InstallOption{
		helm.ValueOverrides(overrides),
		helm.InstallDryRun(dryRun),
		helm.ReleaseName(releaseName),
		helm.InstallReuseName(reuseName),
		helm.InstallDisableHooks(disableHooks),
	}
	return
}

//panicCheck panics if error is not nil
func panicCheck(e error) {
	if e != nil {
		log.Panic(e.Error())
		//panic(e.Error())
	}
}
