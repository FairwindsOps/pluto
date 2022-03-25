// Copyright 2020 Fairwinds
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
// limitations under the License

package cmd

import (
	"fmt"
	"os"

	"github.com/fairwindsops/pluto/v3/pkg/finder"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var (
	directory string
)

func newDetectFilesCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "detect-files",
		Short: "detect-files",
		Long:  `Detect Kubernetes apiVersions in a directory.`,
		Run: func(cmd *cobra.Command, args []string) {

			dir := finder.NewFinder(directory, apiInstance)
			err := dir.FindVersions()
			if err != nil {
				fmt.Println("Error running finder:", err)
				os.Exit(1)
			}

			err = apiInstance.DisplayOutput()
			if err != nil {
				fmt.Println("Error Parsing Output:", err)
				os.Exit(1)
			}
			retCode := apiInstance.GetReturnCode()
			klog.V(5).Infof("retCode: %d", retCode)
			os.Exit(retCode)
		},
	}

	flags := c.Flags()
	flags.StringVarP(&directory, "directory", "d", "", "The directory to scan. If blank, defaults to current working dir.")

	return c
}
