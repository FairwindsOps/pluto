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
	"os"

	"github.com/fairwindsops/pluto/v3/pkg/helm"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var (
	namespace   string
	kubeContext string
)

func newDetectHelmCommand() *cobra.Command {
	c := &cobra.Command{
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

	flags := c.Flags()
	flags.StringVarP(&namespace, "namespace", "n", "", "Only detect releases in a specific namespace.")
	flags.StringVar(&kubeContext, "kube-context", "", "The kube context to use. If blank, defaults to current context.")

	return c
}
