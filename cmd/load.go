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
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

var (
	// loadCmd represents the load command
	loadCmd = &cobra.Command{
		Use:   "load",
		Short: "Load Releases from File to Cluster",
		Long: `This command will decode base64 strings from File into Releases,
	 and 'Helm install' those Releases with the same Chart and Values.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("helm-bulk load called")
			if dryRun {
				log.Println("operating in dry-run mode, helm operations will be run" +
					" in dry-run mode too, meaning installs/upgrades may be returned" +
					" as pending")
			}
			client := helm.NewClient(helm.Host("127.0.0.1:44134"))
			loadedReleases := releases()
			if !nonAuthoritative {
				_, updateReleases := splitReleases(loadedReleases, client)
				purge(updateReleases, client)
			}
			installReleases, updateReleases := splitReleases(loadedReleases, client)
			logReleases(loadedReleases, updateReleases)
			load(installReleases, updateReleases, client)
		},
	}
	dryRun           bool
	nonAuthoritative bool
	fileName         string
)

func init() {
	loadCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false,
		"Uses helm --dry-run internally for any deletes/installs/upgrades")
	loadCmd.Flags().BoolVarP(&nonAuthoritative, "non-authoritive", "n", false,
		"The file that Releases are loaded from is considered non-authoritative,"+
			" so helm-bulk can handle Releases that already exist (using helm upgrade)")
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

//logReleases logs the names of Releases i) loaded from file and (of those)
//ii) currently installed, or a message indicating no Releases were loaded from
//file
func logReleases(releases, existingReleases []*release.Release) {
	if len(releases) > 0 {
		var buffer bytes.Buffer
		buffer.WriteString("Found Helm Releases to load:")
		buffer.WriteString("\n\n")
		addReleasesToBuffer(releases, &buffer)
		if len(existingReleases) > 0 {
			buffer.WriteString("\n")
			buffer.WriteString("..of those, the following are already installed" +
				" (so will result in a 'helm upgrade' rather than 'helm install'):")
			buffer.WriteString("\n\n")
			addReleasesToBuffer(existingReleases, &buffer)
		}
		log.Println(buffer.String())
	} else {
		log.Println("No Helm Releases found to load")
	}
}

//releases decodes the Release file and returns a slice of Releases
func releases() (releases []*release.Release) {
	dat, err := ioutil.ReadFile(fileName)
	panicCheck(err)
	for _, splitString := range strings.Split(string(dat), ",") {
		release, err := utils.DecodeRelease(splitString)
		panicCheck(err)
		releases = append(releases, release)
	}
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
	panicCheck(err)
	for _, release := range releaseResp.GetReleases() {
		existingReleases = append(existingReleases, release)
	}
	panicCheck(err)
	for _, release := range loadedReleases {
		if !containsRelease(release, existingReleases) {
			releases = append(releases, release)
		}
	}
	return
}

//containsRelease returns a bool indicating whether the provided Release is
//in the provided slice
func containsRelease(queryRelease *release.Release,
	targetReleases []*release.Release) (contains bool) {
	for _, release := range targetReleases {
		if queryRelease.GetName() == release.GetName() {
			contains = true
			break
		}
	}
	return
}

//load iterates through first the Releases that need Installing, then those
//that need Upgrading, invoking the func that actually runs through the loading
func load(installReleases, updateReleases []*release.Release,
	client *helm.Client) {
	for _, release := range installReleases {
		loadRelease(release, true, client)
	}
	for _, release := range updateReleases {
		loadRelease(release, false, client)
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
	var buffer bytes.Buffer
	buffer.WriteString("About to purge existing releases:")
	buffer.WriteString("\n\n")
	addReleasesToBuffer(releasesToPurge, &buffer)
	log.Println(buffer.String())
	for _, release := range releasesToPurge {
		releaseName := release.GetName()
		log.Println("Purging Release:", releaseName)
		resp, err := client.DeleteRelease(releaseName, deleteOptions()...)
		panicCheck(err)
		log.Println(releaseName, "helm delete response status:",
			resp.GetRelease().GetInfo().GetStatus().GetCode().String())
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
	deleteOptions = []helm.DeleteOption{
		helm.DeleteDryRun(dryRun),
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
	cv := release.GetConfig()
	overrides := []byte(cv.Raw)
	updateOptions = []helm.UpdateOption{
		helm.UpdateValueOverrides(overrides),
		helm.UpgradeDryRun(dryRun),
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
	releaseName := release.GetName()
	cv := release.GetConfig()
	var overrides = []byte(cv.Raw)
	installOptions = []helm.InstallOption{
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
	}
}
