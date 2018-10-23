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
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// saveCmd represents the save command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show Releases currently stored in the file",
	Long:  `This command will list the Releases currently stored in the file.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("helm-bulk show called")
		show()
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}

//show logs details of Releases it's loaded from file
func show() {
	loadedReleases := Releases()
	var buffer bytes.Buffer
	buffer.WriteString(string(len(loadedReleases)))
	buffer.WriteString(" Releases loaded from file:")
	buffer.WriteString("\n\n")
	for _, release := range loadedReleases {
		buffer.WriteString("    ")
		buffer.WriteString(release.GetName())
		buffer.WriteString("\n")
		buffer.WriteString("        metadata: ")
		buffer.WriteString(strings.Replace(release.GetChart().GetMetadata().String(),
			"\" ", "\"\n                    ", -1))
		buffer.WriteString("\n")
		buffer.WriteString("        values: ")
		buffer.WriteString(strings.Replace(release.GetConfig().String(),
			"\\n", "\"\n                    \"", -1))
		buffer.WriteString("\n\n")
	}
	log.Println(buffer.String())
}
