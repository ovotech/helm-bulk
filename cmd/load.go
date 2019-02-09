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
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mholt/archiver"
	"github.com/ovotech/helm-bulk/utils"
	"github.com/spf13/cobra"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/release"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

var (
	loadCmd = &cobra.Command{
		Use:   "load",
		Short: "Load Releases from File to Cluster",
		Long: `This command will decode base64 strings from File into Releases,
	 and 'Helm install' those Releases with the same Chart and Values.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("helm-bulk load called")
			client := utils.Client(tlsKey, tlsCert, caCert, tlsServerName, disableTLS)
			if dryRun {
				log.Println("*** operating in dry-run mode ***")
			}
			loadedReleases := Releases()
			if len(loadedReleases) > 0 {
				logReleases(loadedReleases, "Helm Releases present in File:")
			} else {
				panic("No Helm Releases found, they're essential for the Load cmd")
			}
			var updateReleases []*release.Release
			if upgrade || delete {
				_, updateReleases = splitReleases(loadedReleases, client)
				if upgrade {
					logReleases(updateReleases, "Existing Helm Releases to update:")
				} else {
					logReleases(updateReleases,
						"Existing Helm Releases to purge (prior to reinstall):")
					purge(updateReleases, client)
				}
			}
			//split Releases a 2nd time, updateReleases may be different now if a delete
			//has just happened.
			installReleases, updateReleases := splitReleases(loadedReleases, client)
			if len(installReleases) > 0 {
				logReleases(installReleases, "Helm Releases to install:")
			} else if !delete && !upgrade {
				log.Println("No Releases found to install, maybe they already exist" +
					" in the Cluster?")
				os.Exit(0)
			} else {
				log.Println("No Releases found to delete or upgrade")
			}
			load(installReleases, updateReleases, client)
		},
	}
	dryRun  bool
	upgrade bool
	delete  bool
)

func init() {
	loadCmd.Flags().BoolVarP(&dryRun, "dry-run", "r", false,
		"Perform a no-op run, essentially just logging to indicate what would be"+
			" done without dry-run enabled")
	loadCmd.Flags().BoolVarP(&upgrade, "upgrade", "u", false,
		"Upgrade existing Releases")
	loadCmd.Flags().BoolVarP(&delete, "delete", "d", false,
		"Delete existing Releases")
	rootCmd.AddCommand(loadCmd)
}

//addRelesesToBuffer adds the string to the buffer, with some added formatting
func addReleasesToBuffer(releases []*release.Release, buffer *bytes.Buffer) *bytes.Buffer {
	for _, release := range releases {
		buffer.WriteString("    ")
		buffer.WriteString(release.GetName())
		buffer.WriteString("\n")
	}
	return buffer
}

func addHeaderToBuffer(header string, buffer *bytes.Buffer) *bytes.Buffer {
	buffer.WriteString(header)
	buffer.WriteString("\n\n")
	return buffer
}

//logReleases logs the names of Releases i) loaded from file and (of those)
//ii) currently installed, or a message indicating no Releases were loaded from
//file
func logReleases(releases []*release.Release, header string) {
	if len(releases) > 0 {
		var buffer bytes.Buffer
		addHeaderToBuffer(header, &buffer)
		addReleasesToBuffer(releases, &buffer)
		log.Println(buffer.String())
	}
}

//Releases decodes the Release file and returns a slice of Releases
func Releases() (releases []*release.Release) {
	wd, err := os.Getwd()
	utils.PanicCheck(err)
	utils.PanicCheck(archiver.TarGz.Open(archiveFilename(), wd))
	dat, err := ioutil.ReadFile(textFilename())
	utils.PanicCheck(err)
	for _, splitString := range strings.Split(string(dat), ",") {
		release, err := utils.DecodeRelease(splitString)
		utils.PanicCheck(err)
		releases = append(releases, release)
	}
	os.Remove(textFilename())
	return
}

//splitReleases obtains a slice of currently installed Releases, which it uses
//along with the provided slice of Releases loaded from file, to compose and
//return two slices; one for Releases to be installed, and another for Releases
//to be upgraded
func splitReleases(loadedReleases []*release.Release,
	client *helm.Client) (releases, existingReleases []*release.Release) {
	var statusFilter = helm.ReleaseListStatuses([]release.Status_Code{
		release.Status_DEPLOYED,
	})
	releaseResp, err := client.ListReleases(statusFilter)
	utils.PanicCheck(err)
	for _, release := range releaseResp.GetReleases() {
		existingReleases = append(existingReleases, release)
	}
	utils.PanicCheck(err)
	for _, release := range loadedReleases {
		if !utils.ContainsRelease(release, existingReleases) {
			releases = append(releases, release)
		}
	}
	return
}

//load iterates through first the Releases that need Installing, then those
//that need Upgrading, invoking the func that actually runs through the loading
func load(installReleases, updateReleases []*release.Release,
	client *helm.Client) {
	if !dryRun {
		for _, release := range installReleases {
			loadRelease(release, true, client)
		}
		for _, release := range updateReleases {
			loadRelease(release, false, client)
		}
	}
}

//loadRelease attempts to Install or Upgrade (depending on whether the Release
//has previously been installed or not) the provided Release.
//If an error is encountered in doing so, it logs the failure and skips to the
//next element in the slice
func loadRelease(release *release.Release, install bool, client *helm.Client) {
	releaseName := release.GetName()
	log.Println("loading Release:", releaseName)
	var statusString string
	var err error
	if install {
		options := installOptions(release)
		var resp *rls.InstallReleaseResponse
		resp, err = client.InstallReleaseFromChart(release.Chart,
			release.GetNamespace(), options...)
		statusString = resp.GetRelease().GetInfo().GetStatus().GetCode().String()
	} else {
		updateOptions := updateOptions(release)
		var resp *rls.UpdateReleaseResponse
		resp, err = client.UpdateReleaseFromChart(releaseName, release.Chart,
			updateOptions...)
		statusString = resp.GetRelease().GetInfo().GetStatus().GetCode().String()
	}
	if err != nil {
		logReleaseFail(releaseName, err)
	} else {
		logReleaseStatusCode(releaseName, statusString, install)
	}
}

//purge deletes the provided releases
func purge(releasesToPurge []*release.Release, client *helm.Client) {
	if !dryRun && len(releasesToPurge) > 0 {
		var buffer bytes.Buffer
		buffer.WriteString("About to purge existing releases:")
		buffer.WriteString("\n\n")
		addReleasesToBuffer(releasesToPurge, &buffer)
		log.Println(buffer.String())
		for _, release := range releasesToPurge {
			releaseName := release.GetName()
			log.Println("Purging Release:", releaseName)
			resp, err := client.DeleteRelease(releaseName, deleteOptions()...)
			utils.PanicCheck(err)
			log.Println(releaseName, "helm delete response status:",
				resp.GetRelease().GetInfo().GetStatus().GetCode().String())
		}
	}
}

//logReleaseFail logs the string and error, with some added formatting
func logReleaseFail(releaseName string, err error) {
	log.Println("loading of Release:", releaseName, " failed:", err.Error())
	log.Println("end of processing for Release:", releaseName)
}

//logReleaseStatusCode logs the strings with some added formatting, including
//what install/upgrade op it relates to
func logReleaseStatusCode(releaseName, statusString string, install bool) {
	var opString string
	if install {
		opString = "install"
	} else {
		opString = "upgrade"
	}
	log.Println(releaseName, "helm", opString, "response status:", statusString)
}

//deleteOptions creates and returns a slice of DeleteOptions
func deleteOptions() (deleteOptions []helm.DeleteOption) {
	deletePurge := true
	deleteDryRun := false
	deleteOptions = []helm.DeleteOption{
		helm.DeleteDryRun(deleteDryRun),
		helm.DeletePurge(deletePurge),
	}
	return
}

//updateOptions creates and returns a slice of UpdateOptions, of which the
//ValueOverrides are obtained from the provided Release
func updateOptions(release *release.Release) (updateOptions []helm.UpdateOption) {
	disableHooks := true
	reuseValues := true
	forceUpgrade := true
	updateDryRun := false
	cv := release.GetConfig()
	overrides := []byte(cv.Raw)
	updateOptions = []helm.UpdateOption{
		helm.UpdateValueOverrides(overrides),
		helm.UpgradeDryRun(updateDryRun),
		helm.ReuseValues(reuseValues),
		helm.UpgradeForce(forceUpgrade),
		helm.UpgradeDisableHooks(disableHooks),
	}
	return
}

//installOptions creates and returns a slice of InstallOptions, of which some
//fields are grabbed from the provided release
func installOptions(release *release.Release) (installOptions []helm.InstallOption) {
	disableHooks := true
	reuseName := true
	installDryRun := false
	releaseName := release.GetName()
	cv := release.GetConfig()
	var overrides = []byte(cv.Raw)
	installOptions = []helm.InstallOption{
		helm.ValueOverrides(overrides),
		helm.InstallDryRun(installDryRun),
		helm.ReleaseName(releaseName),
		helm.InstallReuseName(reuseName),
		helm.InstallDisableHooks(disableHooks),
	}
	return
}
