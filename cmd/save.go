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

	"github.com/ovotech/helm-bulk/utils"
	"github.com/spf13/cobra"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/release"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("helm-bulk save called")
		save()
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// saveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// saveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//save obtains a slice of deployed releases, base64 encodes each release, adds
//the base64 string to a buffer, which it then writes to file.
func save() {
	client := helm.NewClient(helm.Host("127.0.0.1:44134"))
	var statusFilter = helm.ReleaseListStatuses([]release.Status_Code{
		release.Status_DEPLOYED,
	})
	releaseResp, err := client.ListReleases(statusFilter)
	panicCheck(err)
	var buffer bytes.Buffer
	for i, release := range releaseResp.Releases {
		if i > 0 {
			buffer.WriteString(",")
		}
		sEnc, err := utils.EncodeRelease(release)
		panicCheck(err)
		buffer.WriteString(sEnc)
	}
	panicCheck(ioutil.WriteFile("helm-releases.txt", buffer.Bytes(),
		os.FileMode.Perm(0644)))
	panicCheck(ioutil.WriteFile("checksum.txt", md5Hash(buffer.String()),
		os.FileMode.Perm(0644)))
}

//md5Hash returns a byte slice representing the md5 hash of the provided string
func md5Hash(text string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return []byte(hex.EncodeToString(hasher.Sum(nil)))
}
