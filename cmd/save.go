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
var (
	saveCmd = &cobra.Command{
		Use:   "save",
		Short: "Save Releases from Cluster to File",
		Long: `This command will base64 encode current deployed Helm Releases, and
			write them to File.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("helm-bulk save called")
			save()
		},
	}
	orderPrefConfigDir string
)

func init() {
	rootCmd.AddCommand(saveCmd)
	loadCmd.Flags().StringVarP(&orderPrefConfigDir, "order-pref-config-dir", "c", ".",
		"Path (absolute or relative) of directory containing the orderPref.yaml config")
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
	orderPreferences := utils.OrderPref(orderPrefConfigDir)
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
	client := utils.Client(tlsKey, tlsCert, caCert, tlsServerName, disableTLS)
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
