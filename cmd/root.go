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
	"unicode/utf8"

	"github.com/fairwindsops/pluto/pkg/api"
	"github.com/fairwindsops/pluto/pkg/finder"
	"github.com/fairwindsops/pluto/pkg/helm"
	"github.com/rogpeppe/go-internal/semver"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

var (
	version            string
	versionCommit      string
	directory          string
	outputFormat       string
	showAll            bool
	helmVersion        string
	helmStore          string
	ignoreDeprecations bool
	ignoreRemovals     bool
	targetVersion      string
	namespace          string
)

func init() {
	rootCmd.AddCommand(detectFilesCmd)
	rootCmd.PersistentFlags().BoolVarP(&showAll, "show-all", "A", false, "If enabled, will show files that have non-deprecated and non-removed apiVersion. Only applies to tabular output.")
	rootCmd.PersistentFlags().BoolVar(&ignoreDeprecations, "ignore-deprecations", false, "Ignore the default behavior to exit 2 if deprecated apiVersions are found.")
	rootCmd.PersistentFlags().BoolVar(&ignoreRemovals, "ignore-removals", false, "Ignore the default behavior to exit 3 if removed apiVersions are found.")

	rootCmd.PersistentFlags().StringVarP(&targetVersion, "target-version", "t", "v1.16.0", "The version of Kubernetes you wish to check deprecations for.")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "normal", "The output format to use. (normal|wide|json|yaml)")

	detectFilesCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", "The directory to scan. If blank, defaults to current workding dir.")

	rootCmd.AddCommand(detectHelmCmd)
	detectHelmCmd.PersistentFlags().StringVar(&helmVersion, "helm-version", "3", "Helm version in current cluster (2|3)")
	detectHelmCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Only detect releases in a specific namespace.")
	detectHelmCmd.PersistentFlags().StringVar(&helmStore, "helm-store", "configmaps", "Helm storage for v2 (configmaps|secrets)")

	rootCmd.AddCommand(listVersionsCmd)
	rootCmd.AddCommand(detectCmd)

	klog.InitFlags(nil)
	flag.Parse()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

var rootCmd = &cobra.Command{
	Use:   "pluto",
	Short: "pluto",
	Long:  `A tool to detect Kubernetes apiVersions`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You must specify a sub-command.")
		err := cmd.Help()
		if err != nil {
			klog.Error(err)
		}
		os.Exit(1)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		c, _ := utf8.DecodeRuneInString(targetVersion)
		if c != 'v' {
			fmt.Printf("Your --target-version must begin with a 'v'. Got '%s'\n", targetVersion)
			os.Exit(1)
		}

		if !semver.IsValid(targetVersion) {
			fmt.Printf("You must pass a valid semver to --target-version. Got '%s'\n", targetVersion)
			os.Exit(1)
		}
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

		instance := &api.Instance{
			TargetVersion:      targetVersion,
			OutputFormat:       outputFormat,
			ShowAll:            showAll,
			Outputs:            dir.Outputs,
			IgnoreDeprecations: ignoreDeprecations,
			IgnoreRemovals:     ignoreRemovals,
		}

		err = instance.DisplayOutput()
		if err != nil {
			fmt.Println("Error Parsing Output:", err)
			os.Exit(1)
		}
		retCode := instance.GetReturnCode()
		klog.V(5).Infof("retCode: %d", retCode)
		os.Exit(retCode)
	},
}

var detectHelmCmd = &cobra.Command{
	Use:   "detect-helm",
	Short: "detect-helm",
	Long:  `Detect Kubernetes apiVersions in a helm release (in cluster)`,
	Run: func(cmd *cobra.Command, args []string) {
		h := helm.NewHelm(helmVersion, helmStore, namespace)
		err := h.FindVersions()
		if err != nil {
			fmt.Println("Error running helm-detect:", err)
			os.Exit(1)
		}
		instance := &api.Instance{
			TargetVersion:      targetVersion,
			OutputFormat:       outputFormat,
			ShowAll:            showAll,
			IgnoreDeprecations: ignoreDeprecations,
			IgnoreRemovals:     ignoreRemovals,
			Outputs:            h.Outputs,
		}
		err = instance.DisplayOutput()
		if err != nil {
			fmt.Println("Error Parsing Output:", err)
			os.Exit(1)
		}
		retCode := instance.GetReturnCode()
		klog.V(5).Infof("retCode: %d", retCode)
		os.Exit(retCode)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if helmStore != "secrets" && helmStore != "configmaps" && helmVersion == "2" {
			fmt.Println("helm-store should be configmaps or secrets")
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
			instance := &api.Instance{
				TargetVersion:      targetVersion,
				OutputFormat:       outputFormat,
				ShowAll:            showAll,
				IgnoreDeprecations: ignoreDeprecations,
				IgnoreRemovals:     ignoreRemovals,
				Outputs:            output,
			}

			err = instance.DisplayOutput()
			if err != nil {
				fmt.Println("Error parsing output:", err)
				os.Exit(1)
			}
			retCode := instance.GetReturnCode()
			klog.V(5).Infof("retCode: %d", retCode)
			os.Exit(retCode)
		}
		output, err := finder.CheckForAPIVersion(args[0])
		if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}
		instance := &api.Instance{
			TargetVersion:      targetVersion,
			OutputFormat:       outputFormat,
			ShowAll:            showAll,
			IgnoreDeprecations: ignoreDeprecations,
			IgnoreRemovals:     ignoreRemovals,
			Outputs:            output,
		}
		err = instance.DisplayOutput()
		if err != nil {
			fmt.Println("Error parsing output:", err)
			os.Exit(1)
		}
		retCode := instance.GetReturnCode()
		klog.V(5).Infof("retCode: %d", retCode)
		os.Exit(retCode)
	},
}

var listVersionsCmd = &cobra.Command{
	Use:   "list-versions",
	Short: "Outputs a JSON object of the versions that Pluto knows about.",
	Long:  `Outputs a JSON object of the versions that Pluto knows about.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := api.PrintVersionList(outputFormat)
		if err != nil {
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
