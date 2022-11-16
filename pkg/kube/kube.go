// Copyright 2022 FairwindsOps Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	"sync"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	// This is required to auth to cloud providers (i.e. GKE)
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type Kube struct {
	Client kubernetes.Interface
}

var kubeClient *Kube
var once sync.Once

// GetConfigInstance returns a Pluto Kubernetes interface based on the current configuration
func GetConfigInstance(kubeContext string) (*Kube, error) {
	var err error
	var client kubernetes.Interface

	once.Do(func() {
		if kubeClient == nil {
			client, err = GetKubeClient(kubeContext)

			kubeClient = &Kube{
				Client: client,
			}
		}
	})
	if err != nil {
		return nil, err
	}
	return kubeClient, nil
}

// GetKubeClient returns a Kubernetes.Interface based on the current configuration
func GetKubeClient(kubeContext string) (kubernetes.Interface, error) {
	if kubeContext != "" {
		klog.V(3).Infof("using kube context: %s", kubeContext)
	}

	config, err := config.GetConfigWithContext(kubeContext)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// GetKubeDynamicClient returns a dynamic.Interface, rest.Config based on the current configuration
func GetKubeDynamicClient(kubeContext string) (dynamic.Interface, *rest.Config, error) {
	if kubeContext != "" {
		klog.V(3).Infof("using kube context: %s", kubeContext)
	}

	config, err := config.GetConfigWithContext(kubeContext)
	if err != nil {
		return nil, nil, err
	}

	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	return clientset, config, nil
}
