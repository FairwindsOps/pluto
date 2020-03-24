package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

var (
	// Main
	version       string
	versionCommit string
)

func init() {
	klog.InitFlags(nil)
	flag.Parse()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "run",
	Long:  `A tool to detect Kubernetes apiVersions in a directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello world")
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
