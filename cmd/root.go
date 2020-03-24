package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/fairwindsops/api-version-finder/pkg/finder"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

var (
	version       string
	versionCommit string
	directory     string
)

func init() {
	rootCmd.AddCommand(detectCmd)
	detectCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", "The directory to scan. If blank, defaults to current workding dir.")
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
		fmt.Println("Thank you for using the tool. Running....")
		if directory == "" {
			directory, err := os.Getwd()
			if err != nil {
				klog.Fatal(err)
			}
			klog.V(3).Infof("using working directory: %s", directory)
		}
		err := finder.FindFiles(directory)
		if err != nil {
			fmt.Printf("There was an error: %s", err)
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
