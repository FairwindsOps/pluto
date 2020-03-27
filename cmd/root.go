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
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fairwindsops/pluto/pkg/api"
	"github.com/fairwindsops/pluto/pkg/finder"
	"github.com/fairwindsops/pluto/pkg/helm"

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
	helmVersion       string
)

func init() {
	rootCmd.AddCommand(detectFilesCmd)
	rootCmd.PersistentFlags().BoolVar(&showNonDeprecated, "show-non-deprecated", false, "If enabled, will show files that have non-deprecated apiVersion. Only applies to tabular output.")

	detectFilesCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", "The directory to scan. If blank, defaults to current workding dir.")
	detectFilesCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "tabular", "The output format to use. (tabular|json|yaml)")

	rootCmd.AddCommand(detectHelmCmd)
	detectHelmCmd.PersistentFlags().StringVar(&helmVersion, "helm-version", "3", "Helm version in current cluster (2|3)")

	rootCmd.AddCommand(detectCmd)

	klog.InitFlags(nil)
	flag.Parse()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

// Checker is an interface to find versions
type Checker interface {
	FindVersions() error
}

var rootCmd = &cobra.Command{
	Use:   "pluto",
	Short: "pluto",
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

var detectFilesCmd = &cobra.Command{
	Use:   "detect-files",
	Short: "detect-files",
	Long:  `Detect Kubernetes apiVersions in a directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := finder.NewFinder(directory)
		err := dir.FindVersions()
		if err != nil {
			fmt.Println("Error running finder:", err)
			os.Exit(1)
		}

		if dir.Outputs == nil {
			fmt.Println("No api-versioned files found in specified directory.")
			os.Exit(0)
		}

		err = api.DisplayOutput(dir.Outputs, outputFormat, showNonDeprecated)
		if err != nil {
			fmt.Println("Error Parsing Output:", err)
			os.Exit(1)
		}
	},
}

var detectHelmCmd = &cobra.Command{
	Use:   "detect-helm",
	Short: "detect-helm",
	Long:  `Detect Kubernetes apiVersions in a helm release (in cluster)`,
	Run: func(cmd *cobra.Command, args []string) {
		h := helm.NewHelm(helmVersion)
		err := h.FindVersions()
		if err != nil {
			fmt.Println("Error running helm-detect:", err)
			os.Exit(1)
		}
		err = api.DisplayOutput(h.Outputs, outputFormat, showNonDeprecated)
		if err != nil {
			fmt.Println("Error Parsing Output:", err)
			os.Exit(1)
		}
	},
}

var detectCmd = &cobra.Command{
	Use:   "detect [file to check or -]",
	Short: "Checks a single file or stdin for deprecated apiVersions.",
	Long:  `Detect deprecated apiVersion in a specific file or other input. Accepts multi-document yaml files and/or - for stdin. Useful for helm testing.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires a file argument")
		}
		if isFileOrStdin(args[0]) {
			return nil
		}
		return fmt.Errorf("invalid file specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		klog.V(3).Infof("arg0: %s", args[0])

		if args[0] == "-" {
			//stdin
			fileData, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println("Error reading stdin:", err)
				os.Exit(1)
			}
			output, err := api.IsVersioned(fileData)
			if err != nil {
				fmt.Println("Error checking for versions:", err)
				os.Exit(1)
			}
			err = api.DisplayOutput(output, outputFormat, showNonDeprecated)
			if err != nil {
				fmt.Println("Error parsing output:", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
		output, err := finder.CheckForAPIVersion(args[0])
		if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}
		err = api.DisplayOutput(output, outputFormat, showNonDeprecated)
		if err != nil {
			fmt.Println("Error parsing output:", err)
			os.Exit(1)
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

func isFileOrStdin(name string) bool {
	if name == "-" {
		return true
	}
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
