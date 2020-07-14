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
	"strings"

	"github.com/fairwindsops/pluto/pkg/api"
	"github.com/fairwindsops/pluto/pkg/finder"
	"github.com/fairwindsops/pluto/pkg/helm"
	"github.com/rogpeppe/go-internal/semver"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

var (
	version                string
	versionCommit          string
	additionalVersionsFile string
	directory              string
	outputFormat           string
	helmVersion            string
	helmStore              string
	ignoreDeprecations     bool
	ignoreRemovals         bool
	namespace              string
	apiInstance            *api.Instance
	targetVersions         map[string]string
	customColumns          []string
)

var outputOptions = []string{
	"json",
	"yaml",
	"normal",
	"wide",
	"custom",
}

var possibleColumns = []string{
	"NAME",
	"NAMESPACE",
	"KIND",
	"VERSION",
	"REPLACEMENT",
	"DEPRECATED",
	"DEPRECATED IN",
	"REMOVED",
	"REMOVED IN",
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&ignoreDeprecations, "ignore-deprecations", false, "Ignore the default behavior to exit 2 if deprecated apiVersions are found.")
	rootCmd.PersistentFlags().BoolVar(&ignoreRemovals, "ignore-removals", false, "Ignore the default behavior to exit 3 if removed apiVersions are found.")
	rootCmd.PersistentFlags().StringVarP(&additionalVersionsFile, "additional-versions", "f", "", "Additional deprecated versions file to add to the list. Cannot contain any existing versions")
	rootCmd.PersistentFlags().StringToStringVarP(&targetVersions, "target-versions", "t", targetVersions, "A map of targetVersions to use. This flag supersedes all defaults in version files.")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "normal", "The output format to use. (normal|wide|custom|json|yaml)")
	rootCmd.PersistentFlags().StringSliceVar(&customColumns, "columns", nil, "A list of columns to print when using --output custom")

	rootCmd.AddCommand(detectFilesCmd)
	detectFilesCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", "The directory to scan. If blank, defaults to current working dir.")

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
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		//verify output option
		if !stringInSlice(outputFormat, outputOptions) {
			return fmt.Errorf("--output must be one of %v", outputOptions)
		}

		if outputFormat == "custom" {
			if len(customColumns) < 1 {
				return fmt.Errorf("when --output=custom you must specify --columns")
			}
			// Uppercase all columns entered on CLI
			var tempColumns []string
			for _, colString := range customColumns {
				tempColumns = append(tempColumns, strings.ToUpper(colString))
			}

			customColumns = tempColumns
			for _, c := range customColumns {
				if !stringInSlice(c, possibleColumns) {
					return fmt.Errorf("invalid custom column option %s - must be one of %v", c, possibleColumns)
				}
			}
		}

		defaultVersions, defaultTargetVersions, err := api.GetDefaultVersionList()
		if err != nil {
			return err
		}
		var deprecatedVersionList []api.Version
		if additionalVersionsFile != "" {
			klog.V(2).Infof("looking for versions file: %s", additionalVersionsFile)
			data, err := ioutil.ReadFile(additionalVersionsFile)
			if err != nil {
				return err
			}
			additionalVersions, additionalTargetVersions, err := api.UnMarshalVersions(data)
			if err != nil {
				return err
			}
			deprecatedVersionList, err = api.CombineAdditionalVersions(additionalVersions, defaultVersions)
			if err != nil {
				return err
			}
			for c, v := range additionalTargetVersions {
				klog.V(2).Infof("received target version from config: %s %s", c, v)
				// Only add them to default target versions if they do not supersed any previously included
				// This prevents overwriting the internal defaults
				if _, found := defaultTargetVersions[c]; !found {
					defaultTargetVersions[c] = v
				}
			}
		} else {
			klog.V(2).Info("no additional versions needed")
			deprecatedVersionList = defaultVersions
		}

		// From this compiled list of deprecations, build a component list
		var componentList []string
		for _, v := range deprecatedVersionList {
			if !stringInSlice(v.Component, componentList) {
				componentList = append(componentList, v.Component)
			}
		}

		// Combine the default target versions and the ones that are passed. Ones that are passed in take precedence
		if targetVersions != nil {
			for k, v := range defaultTargetVersions {
				if _, found := targetVersions[k]; !found {
					klog.V(2).Infof("assuming default targetVersion %s %s", k, v)
					targetVersions[k] = v
				}
			}
		} else {
			targetVersions = defaultTargetVersions
		}

		// Verify that we have valid target versions for all components
		for component, version := range targetVersions {
			if !semver.IsValid(version) {
				return fmt.Errorf("you must use valid semver for all target versions with a leading 'v' - got %s %s", component, version)
			}
		}
		for _, c := range componentList {
			if _, found := targetVersions[c]; !found {
				return fmt.Errorf("you must pass a targetVersion for every component in the list - missing component: %s", c)
			}
		}

		apiInstance = &api.Instance{
			TargetVersions:     targetVersions,
			OutputFormat:       outputFormat,
			CustomColumns:      customColumns,
			IgnoreDeprecations: ignoreDeprecations,
			IgnoreRemovals:     ignoreRemovals,
			DeprecatedVersions: deprecatedVersionList,
		}

		return nil
	},
}

var detectFilesCmd = &cobra.Command{
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

var detectHelmCmd = &cobra.Command{
	Use:   "detect-helm",
	Short: "detect-helm",
	Long:  `Detect Kubernetes apiVersions in a helm release (in cluster)`,
	Run: func(cmd *cobra.Command, args []string) {
		h := helm.NewHelm(helmVersion, helmStore, namespace, apiInstance)
		err := h.FindVersions()
		if err != nil {
			fmt.Println("Error running helm-detect:", err)
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

var listVersionsCmd = &cobra.Command{
	Use:   "list-versions",
	Short: "Outputs a JSON object of the versions that Pluto knows about.",
	Long:  `Outputs a JSON object of the versions that Pluto knows about.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := apiInstance.PrintVersionList(outputFormat)
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

func stringInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if s == v {
			return true
		}
	}
	return false
}
