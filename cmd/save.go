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
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/mholt/archiver"
	"github.com/ovotech/helm-bulk/utils"
	"github.com/spf13/cobra"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/release"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save Releases from Cluster to File",
	Long: `This command will base64 encode current deployed Helm Releases, and
			write them to File.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("helm-bulk save called")
		save()
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
}

//releaseFromName returns the Release in the provided slice for which the Name
//matches the provided searchName string. If none match, nil is returned.
func releaseFromName(searchName string,
	releases []*release.Release) (targetRelease *release.Release) {
	for _, release := range releases {
		if release.Name == searchName {
			targetRelease = release
			break
		}
	}
	return
}

//targetReleases returns a slice of Releases based on the preferred ordering and
//those currently installed.
func targetReleases(releases []*release.Release) (targetReleases []*release.Release) {
	orderPreferences := utils.OrderPref()
	for _, orderedReleaseName := range orderPreferences {
		targetRelease := releaseFromName(orderedReleaseName, releases)
		if targetRelease != nil {
			targetReleases = append(targetReleases, targetRelease)
		}
	}
	for _, release := range releases {
		if !utils.ContainsRelease(release, targetReleases) {
			targetReleases = append(targetReleases, release)
		}
	}
	return
}

//save obtains a slice of deployed releases, base64 encodes each release, adds
//the base64 string to a buffer, which it then writes to file.
func save() {
	client := utils.Client()
	var statusFilter = helm.ReleaseListStatuses([]release.Status_Code{
		release.Status_DEPLOYED,
	})
	releaseResp, err := client.ListReleases(statusFilter)
	utils.PanicCheck(err)
	var buffer bytes.Buffer
	releases := releaseResp.GetReleases()
	targetReleases := targetReleases(releases)
	for i, release := range targetReleases {
		if i > 0 {
			buffer.WriteString(",")
		}
		sEnc, errb := utils.EncodeRelease(release)
		utils.PanicCheck(errb)
		buffer.WriteString(sEnc)
	}
	utils.PanicCheck(ioutil.WriteFile(textFilename(), buffer.Bytes(),
		os.FileMode.Perm(0644)))
	utils.PanicCheck(archiver.TarGz.Make(archiveFilename(),
		[]string{textFilename()}))
	log.Println("Wrote " + strconv.Itoa(len(releases)) + " Helm Releases to file")
	os.Remove(textFilename())
}

//md5Hash returns a byte slice representing the md5 hash of the provided string
func md5Hash(text string) (hash string) {
	hasher := md5.New()
	hasher.Write([]byte(text))
	hash = hex.EncodeToString(hasher.Sum(nil))
	return
}
