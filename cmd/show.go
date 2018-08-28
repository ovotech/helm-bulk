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
