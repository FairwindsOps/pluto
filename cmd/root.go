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
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"gopkg.in/yaml.v3"

	"github.com/fairwindsops/api-version-finder/pkg/finder"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

var (
	version           string
	versionCommit     string
	directory         string
	outputFormat      string
	showNonDeprecated bool
)

func init() {
	rootCmd.AddCommand(detectCmd)
	detectCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", "The directory to scan. If blank, defaults to current workding dir.")
	detectCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "tabular", "The output format to use. (tabular|json|yaml)")
	detectCmd.PersistentFlags().BoolVar(&showNonDeprecated, "show-non-deprecated", false, "If enabled, will show files that have non-deprecated apiVersion. Only applies to tabular output.")

	klog.InitFlags(nil)
	flag.Parse()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

var rootCmd = &cobra.Command{
	Use:   "api-version-finder",
	Short: "api-version-finder",
	Long:  `A tool to detect Kubernetes apiVersions`,
	Run: func(cmd *cobra.Command, args []string) {
		klog.Error("You must specify a sub-command.")
		err := cmd.Help()
		if err != nil {
			klog.Error(err)
		}
		os.Exit(1)
	},
}

var detectCmd = &cobra.Command{
	Use:   "detect-files",
	Short: "detect-files",
	Long:  `Detect Kubernetes apiVersions in a directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := finder.NewFinder(directory)
		err := dir.FindVersions()
		if err != nil {
			fmt.Printf("Error running finder: %s\n", err.Error())
			os.Exit(1)
		}

		if dir.APIFiles == nil {
			fmt.Println("No api-versioned files found in specified directory.")
			os.Exit(0)
		}

		switch outputFormat {
		case "tabular":
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintln(w, "KIND\t VERSION\t DEPRECATED\t FILE")

			for _, file := range dir.APIFiles {
				// Don't show non-deprecated apis if we have them disabled
				if !showNonDeprecated {
					if !file.APIVersion.Deprecated {
						continue
					}
				}
				kind := file.APIVersion.Kind
				deprecated := fmt.Sprintf("%t", file.APIVersion.Deprecated)
				version := file.APIVersion.Name
				fileName := file.Name

				fmt.Fprintf(w, "%s\t %s\t %s\t %s\t\n", kind, version, deprecated, fileName)
			}
			w.Flush()
		case "json":
			outData, err := json.Marshal(dir.APIFiles)
			if err != nil {
				fmt.Printf("Error generating JSON: %s\n", err.Error())
				os.Exit(1)
			}
			fmt.Println(string(outData))
		case "yaml":
			outData, err := yaml.Marshal(dir.APIFiles)
			if err != nil {
				fmt.Printf("Error generating yaml: %s\n", err.Error())
				os.Exit(1)
			}
			fmt.Println(string(outData))
		}
	},
}

// Execute the stuff
func Execute(VERSION string, COMMIT string) {
	version = VERSION
	versionCommit = COMMIT
	if err := rootCmd.Execute(); err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}
