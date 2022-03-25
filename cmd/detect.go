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
	"io/ioutil"
	"os"

	"github.com/fairwindsops/pluto/v3/pkg/api"
	"github.com/fairwindsops/pluto/v3/pkg/finder"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

func newDetectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "detect [file to check or -]",
		Short: "Checks a single file or stdin for deprecated apiVersions.",
		Long:  `Detect deprecated apiVersion in a specific file or other input. Accepts multi-document yaml files and/or - for stdin. Useful for helm testing.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires a file argument")
			}
			if api.IsFileOrStdin(args[0]) {
				return nil
			}
			return fmt.Errorf("invalid file specified: %s", args[0])
		},
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			klog.V(3).Infof("arg0: %s", args[0])

			if args[0] == "-" {
				//stdin
				fileData, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					fmt.Println("Error reading stdin:", err)
					os.Exit(1)
				}

				apiInstance.Outputs, err = apiInstance.IsVersioned(fileData)
				if err != nil {
					fmt.Println("Error checking for versions:", err)
					os.Exit(1)
				}

				err = apiInstance.DisplayOutput()
				if err != nil {
					fmt.Println("Error parsing output:", err)
					os.Exit(1)
				}
				retCode := apiInstance.GetReturnCode()
				klog.V(5).Infof("retCode: %d", retCode)
				os.Exit(retCode)
			}

			// File input
			dir := finder.Dir{
				Instance: apiInstance,
			}
			apiInstance.Outputs, err = dir.CheckForAPIVersion(args[0])
			if err != nil {
				fmt.Println("Error reading file:", err)
				os.Exit(1)
			}

			err = apiInstance.DisplayOutput()
			if err != nil {
				fmt.Println("Error parsing output:", err)
				os.Exit(1)
			}
			retCode := apiInstance.GetReturnCode()
			klog.V(5).Infof("retCode: %d", retCode)
			os.Exit(retCode)
		},
	}
}
