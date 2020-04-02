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

package helm

import (
	"fmt"
	"os"
	"sync"

	"k8s.io/client-go/kubernetes"
	// This is required to auth to gcp (i.e. GKE)
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type kube struct {
	Client kubernetes.Interface
}

var kubeClient *kube
var once sync.Once

// GetConfigInstance returns a Kubernetes interface based on the current configuration
func getConfigInstance() *kube {
	once.Do(func() {
		if kubeClient == nil {
			kubeClient = &kube{
				Client: getKubeClient(),
			}
		}
	})
	return kubeClient
}

func getKubeClient() kubernetes.Interface {
	kubeConf, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error getting kubeconfig:", err)
		os.Exit(1)
	}
	clientset, err := kubernetes.NewForConfig(kubeConf)
	if err != nil {
		fmt.Println("Error creating kubernetes client:", err)
		os.Exit(1)
	}
	return clientset
}
