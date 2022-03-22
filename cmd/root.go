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

	"github.com/fairwindsops/pluto/v3/pkg/api"
	"github.com/fairwindsops/pluto/v3/pkg/finder"
	"github.com/fairwindsops/pluto/v3/pkg/helm"
	"github.com/rogpeppe/go-internal/semver"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

var (
	version                string
	versionCommit          string
	versionFileData        []byte
	additionalVersionsFile string
	directory              string
	outputFormat           string
	ignoreDeprecations     bool
	ignoreRemovals         bool
	namespace              string
	apiInstance            *api.Instance
	targetVersions         map[string]string
	customColumns          []string
	componentsFromUser     []string
	onlyShowRemoved        bool
	kubeContext            string
)

const (
	envPrefix = "PLUTO"
)

var outputOptions = []string{
	"json",
	"yaml",
	"normal",
	"wide",
	"custom",
	"markdown",
	"csv",
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&ignoreDeprecations, "ignore-deprecations", false, "Ignore the default behavior to exit 2 if deprecated apiVersions are found.")
	rootCmd.PersistentFlags().BoolVar(&ignoreRemovals, "ignore-removals", false, "Ignore the default behavior to exit 3 if removed apiVersions are found.")
	rootCmd.PersistentFlags().BoolVarP(&onlyShowRemoved, "only-show-removed", "r", false, "Only display the apiVersions that have been removed in the target version.")
	rootCmd.PersistentFlags().StringVarP(&additionalVersionsFile, "additional-versions", "f", "", "Additional deprecated versions file to add to the list. Cannot contain any existing versions")
	rootCmd.PersistentFlags().StringToStringVarP(&targetVersions, "target-versions", "t", targetVersions, "A map of targetVersions to use. This flag supersedes all defaults in version files.")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "normal", "The output format to use. (normal|wide|custom|json|yaml|markdown|csv)")
	rootCmd.PersistentFlags().StringSliceVar(&customColumns, "columns", nil, "A list of columns to print. Mandatory when using --output custom, optional with --output markdown")
	rootCmd.PersistentFlags().StringSliceVar(&componentsFromUser, "components", nil, "A list of components to run checks for. If nil, will check for all found in versions.")

	rootCmd.AddCommand(detectFilesCmd)
	detectFilesCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", "The directory to scan. If blank, defaults to current working dir.")

	rootCmd.AddCommand(detectHelmCmd)
	detectHelmCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Only detect releases in a specific namespace.")
	detectHelmCmd.PersistentFlags().StringVar(&kubeContext, "kube-context", "", "The kube context to use. If blank, defaults to current context.")

	rootCmd.AddCommand(listVersionsCmd)
	rootCmd.AddCommand(detectCmd)

	klog.InitFlags(nil)
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("v"))
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	bindFlags(cmd, v)

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			err := v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
			if err != nil {
				klog.Errorf("error binding flag %s to env var %s_%s: %w", f.Name, envPrefix, envVarSuffix, err)
				return
			}
		}

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				klog.Errorf("error setting flag %s to %v: %w", f.Name, val, err)
				return
			}
		}
	})
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
		err := initializeConfig(cmd)
		if err != nil {
			return err
		}

		//verify output option
		if !api.StringInSlice(outputFormat, outputOptions) {
			return fmt.Errorf("--output must be one of %v", outputOptions)
		}

		if outputFormat == "custom" {
			if len(customColumns) < 1 {
				return fmt.Errorf("when --output=custom you must specify --columns")
			}
		}

		if outputFormat == "custom" || (outputFormat == "markdown" && len(customColumns) >= 1) {
			// Uppercase all columns entered on CLI
			var tempColumns []string
			for _, colString := range customColumns {
				tempColumns = append(tempColumns, strings.ToUpper(colString))
			}

			customColumns = tempColumns
			for _, c := range customColumns {
				if !api.StringInSlice(c, api.PossibleColumnNames) {
					return fmt.Errorf("invalid custom column option %s - must be one of %v", c, api.PossibleColumnNames)
				}
			}
		}

		defaultVersions, defaultTargetVersions, err := api.GetDefaultVersionList(versionFileData)
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

		// From the compiled list of deprecations and the components flag, build a component list
		var componentList []string
		for _, v := range deprecatedVersionList {
			if !api.StringInSlice(v.Component, componentList) {
				if componentsFromUser == nil {
					// if the pass-in components are nil, then we use the versions in the main list
					componentList = append(componentList, v.Component)
				} else {
					if api.StringInSlice(v.Component, componentsFromUser) {
						// if a component list was passed, make sure that it is in the
						// the list of versions before adding it to the final list
						componentList = append(componentList, v.Component)
					}
				}
			}
		}
		if len(componentList) < 1 {
			return fmt.Errorf("cannot find deprecations for zero components")
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

		// verify that we have valid target versions for all components
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

		// this apiInstance will be used by all detection methods
		apiInstance = &api.Instance{
			TargetVersions:     targetVersions,
			OutputFormat:       outputFormat,
			CustomColumns:      customColumns,
			IgnoreDeprecations: ignoreDeprecations,
			IgnoreRemovals:     ignoreRemovals,
			OnlyShowRemoved:    onlyShowRemoved,
			DeprecatedVersions: deprecatedVersionList,
			Components:         componentList,
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
		h, err := helm.NewHelm(namespace, kubeContext, apiInstance)
		if err != nil {
			fmt.Printf("error getting helm configuration: %s\n", err.Error())
			os.Exit(1)
		}
		err = h.FindVersions()
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
}

var detectCmd = &cobra.Command{
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
func Execute(VERSION string, COMMIT string, versionsFile []byte) {
	version = VERSION
	versionCommit = COMMIT
	versionFileData = versionsFile
	if err := rootCmd.Execute(); err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}
