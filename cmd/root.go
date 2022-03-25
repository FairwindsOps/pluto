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
	"strings"

	"github.com/fairwindsops/pluto/v3/pkg/api"
	"github.com/rogpeppe/go-internal/semver"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

var (
	versionFileData        []byte
	additionalVersionsFile string
	outputFormat           string
	ignoreDeprecations     bool
	ignoreRemovals         bool
	apiInstance            *api.Instance
	targetVersions         map[string]string
	customColumns          []string
	componentsFromUser     []string
	onlyShowRemoved        bool
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

func NewRootCommand(VERSION string, COMMIT string, versionsFile []byte) (*cobra.Command, error) {
	version := VERSION
	versionCommit := COMMIT
	versionFileData = versionsFile

	c := &cobra.Command{
		Use:               "pluto",
		Short:             "pluto",
		Long:              `A tool to detect Kubernetes apiVersions`,
		PersistentPreRunE: preRunCheck,
	}

	persistentFlags := c.PersistentFlags()
	persistentFlags.BoolVar(&ignoreDeprecations, "ignore-deprecations", false, "Ignore the default behavior to exit 2 if deprecated apiVersions are found.")
	persistentFlags.BoolVar(&ignoreRemovals, "ignore-removals", false, "Ignore the default behavior to exit 3 if removed apiVersions are found.")
	persistentFlags.BoolVarP(&onlyShowRemoved, "only-show-removed", "r", false, "Only display the apiVersions that have been removed in the target version.")
	persistentFlags.StringVarP(&additionalVersionsFile, "additional-versions", "f", "", "Additional deprecated versions file to add to the list. Cannot contain any existing versions")
	persistentFlags.StringToStringVarP(&targetVersions, "target-versions", "t", targetVersions, "A map of targetVersions to use. This flag supersedes all defaults in version files.")
	persistentFlags.StringVarP(&outputFormat, "output", "o", "normal", "The output format to use. (normal|wide|custom|json|yaml|markdown|csv)")
	persistentFlags.StringSliceVar(&customColumns, "columns", nil, "A list of columns to print. Mandatory when using --output custom, optional with --output markdown")
	persistentFlags.StringSliceVar(&componentsFromUser, "components", nil, "A list of components to run checks for. If nil, will check for all found in versions.")

	c.AddCommand(
		newDetectCommand(),
		newDetectFilesCommand(),
		newDetectHelmCommand(),
		newListVersionsCommand(),
		newVersionsCommand(version, versionCommit),
	)

	klog.InitFlags(nil)
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("v"))

	return c, nil
}

func preRunCheck(cmd *cobra.Command, args []string) error {
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
