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

	"github.com/fairwindsops/api-version-finder/pkg/api"
	"github.com/fairwindsops/api-version-finder/pkg/finder"
	"github.com/fairwindsops/api-version-finder/pkg/helm"

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
	detectFilesCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", "The directory to scan. If blank, defaults to current workding dir.")
	detectFilesCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "tabular", "The output format to use. (tabular|json|yaml)")
	detectFilesCmd.PersistentFlags().BoolVar(&showNonDeprecated, "show-non-deprecated", false, "If enabled, will show files that have non-deprecated apiVersion. Only applies to tabular output.")

	rootCmd.AddCommand(detectHelmCmd)
	detectHelmCmd.PersistentFlags().StringVar(&helmVersion, "helm-version", "3", "Helm version in current cluster (2|3)")
	detectHelmCmd.PersistentFlags().BoolVar(&showNonDeprecated, "show-non-deprecated", false, "If enabled, will show files that have non-deprecated apiVersion. Only applies to tabular output.")

	klog.InitFlags(nil)
	flag.Parse()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

type Checker interface {
	FindVersions() error
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

var detectFilesCmd = &cobra.Command{
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

		if dir.Outputs == nil {
			fmt.Println("No api-versioned files found in specified directory.")
			os.Exit(0)
		}

		err = parseOutput(dir.Outputs)
		if err != nil {
			klog.Errorf("Error Parsing Output: %v", err)
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
			klog.Fatalf("ERROR: %v", err)
		}
		err = parseOutput(h.Outputs)
		if err != nil {
			klog.Errorf("Error Parsing Output: %v", err)
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

func parseOutput(outputs []*api.Output) error {
	var err error
	var outData []byte
	switch outputFormat {
	case "tabular":
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 2, ' ', 0)
		_, err = fmt.Fprintln(w, "KIND\t VERSION\t DEPRECATED\t RESOURCE NAME")

		for _, output := range outputs {
			// Don't show non-deprecated apis if we have them disabled
			if !showNonDeprecated {
				if !output.APIVersion.Deprecated {
					continue
				}
			}
			kind := output.APIVersion.Kind
			deprecated := fmt.Sprintf("%t", output.APIVersion.Deprecated)
			version := output.APIVersion.Name
			fileName := output.Name

			_, err = fmt.Fprintf(w, "%s\t %s\t %s\t %s\t\n", kind, version, deprecated, fileName)
		}
		err = w.Flush()
	case "json":
		outData, err = json.Marshal(outputs)
	case "yaml":
		outData, err = yaml.Marshal(outputs)
	}
	if err != nil {
		return err
	}
	fmt.Println(string(outData))
	return nil
}
