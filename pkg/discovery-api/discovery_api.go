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

package discoveryapi

import (
	"context"
	"encoding/json"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/fairwindsops/pluto/v5/pkg/api"
)

// DiscoveryClient is the declaration to hold objects needed for client-go/discovery.
type DiscoveryClient struct {
	ClientSet       dynamic.Interface
	restConfig      *rest.Config
	DiscoveryClient discovery.DiscoveryInterface
	Instance        *api.Instance
}

// NewDiscoveryClient returns a new struct with config portions complete.
func NewDiscoveryClient(instance *api.Instance) (*DiscoveryClient, error) {
	cl := &DiscoveryClient{
		Instance: instance,
	}

	var err error
	cl.restConfig, err = NewRestClientConfig(rest.InClusterConfig)
	if err != nil {
		return nil, err
	}

	if cl.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(cl.restConfig); err != nil {
		return nil, err
	}

	cl.ClientSet, err = dynamic.NewForConfig(cl.restConfig)
	if err != nil {
		return nil, err
	}
	return cl, nil
}

// NewRestClientConfig returns a new Rest Client config portions complete.
func NewRestClientConfig(inClusterFn func() (*rest.Config, error)) (*rest.Config, error) {

	if restConfig, err := inClusterFn(); err == nil {
		return restConfig, nil
	}

	pathOptions := clientcmd.NewDefaultPathOptions()

	config, err := pathOptions.GetStartingConfig()
	if err != nil {
		return nil, err
	}

	configOverrides := clientcmd.ConfigOverrides{}

	clientConfig := clientcmd.NewDefaultClientConfig(*config, &configOverrides)
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return restConfig, nil
}

// GetApiResources discovers the api-resources for a cluster
func (cl *DiscoveryClient) GetApiResources() error {

	resourcelist, err := cl.DiscoveryClient.ServerPreferredResources()
	if err != nil {
		if apierrors.IsNotFound(err) {
			return err
		}
		if apierrors.IsForbidden(err) {
			klog.Error("Failed to list objects for Name discovery. Permission denied! Please check if you have the proper authorization")
			return err
		}
	}

	gvrs := []schema.GroupVersionResource{}
	for _, rl := range resourcelist {
		for i := range rl.APIResources {
			gv, _ := schema.ParseGroupVersion(rl.GroupVersion)
			ResourceName := rl.APIResources[i].Name
			g := schema.GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: ResourceName}
			gvrs = append(gvrs, g)
		}
	}

	var results []map[string]interface{}
	for _, g := range gvrs {
		ri := cl.ClientSet.Resource(g)
		klog.V(2).Infof("Retrieving : %s.%s.%s", g.Resource, g.Version, g.Group)
		rs, err := ri.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			klog.Error("Failed to retrieve: ", g, err)
			continue
		}

		if len(rs.Items) == 0 {
			klog.V(2).Infof("No annotations for ResourceVer %s", rs.GetAPIVersion())
			obj := rs.UnstructuredContent()
			data, err := json.Marshal(obj)
			if err != nil {
				klog.Error("Failed to marshal data ", err.Error())
				return err
			}
			output, err := cl.Instance.IsVersioned(data)
			if err != nil {
				return err
			}
			if output == nil {
				continue
			}
			cl.Instance.Outputs = append(cl.Instance.Outputs, output...)

		} else {
			for _, r := range rs.Items {
				if jsonManifest, ok := r.GetAnnotations()["kubectl.kubernetes.io/last-applied-configuration"]; ok {
					var manifest map[string]interface{}

					err := json.Unmarshal([]byte(jsonManifest), &manifest)
					if err != nil {
						klog.Error("failed to parse 'last-applied-configuration' annotation of resource %s/%s: %s", r.GetNamespace(), r.GetName(), err.Error())
						continue
					}
					data, err := json.Marshal(manifest)
					if err != nil {
						klog.Error("Failed to marshal data ", err.Error())
						return err
					}
					output, err := cl.Instance.IsVersioned(data)
					if err != nil {
						return err
					}
					cl.Instance.Outputs = append(cl.Instance.Outputs, output...)
				}
			}
		}

	}

	klog.V(6).Infof("Result from resources: %d", len(results))
	return nil
}
