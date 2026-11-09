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
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	discoveryFake "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/dynamic/fake"

	"github.com/fairwindsops/pluto/v5/pkg/api"
)

// preferredResourcesDiscovery returns a fixed resource list, because the
// fake discovery client always returns an empty one.
type preferredResourcesDiscovery struct {
	discovery.DiscoveryInterface
	resources []*metav1.APIResourceList
}

func (d *preferredResourcesDiscovery) ServerPreferredResources() ([]*metav1.APIResourceList, error) {
	return d.resources, nil
}

func TestNewDiscoveryAPIClientValidEmpty(t *testing.T) {

	scheme := runtime.NewScheme()
	clientset := fake.NewSimpleDynamicClient(scheme)
	discoveryClient := discoveryFake.FakeDiscovery{}
	testOpts := DiscoveryClient{
		ClientSet:       clientset,
		DiscoveryClient: &discoveryClient,
	}

	err := testOpts.GetApiResources()
	if err != nil {
		t.Errorf("Unable to fetch API Resources")
	}

}

func TestGetApiResourcesListsOnlyKindsWithRules(t *testing.T) {
	deployments := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	roles := schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"}
	configMaps := schema.GroupVersionResource{Version: "v1", Resource: "configmaps"}

	deployment := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]any{
			"name":      "web",
			"namespace": "default",
			"annotations": map[string]any{
				"kubectl.kubernetes.io/last-applied-configuration": `{"apiVersion":"extensions/v1beta1","kind":"Deployment","metadata":{"name":"web","namespace":"default"}}`,
			},
		},
	}}

	clientset := fake.NewSimpleDynamicClientWithCustomListKinds(runtime.NewScheme(), map[schema.GroupVersionResource]string{
		deployments: "DeploymentList",
		roles:       "RoleList",
		configMaps:  "ConfigMapList",
	}, deployment)

	extensionsDeployment := api.Version{Name: "extensions/v1beta1", Kind: "Deployment"}
	instance := &api.Instance{DeprecatedVersions: []api.Version{
		extensionsDeployment,
		{Name: "rbac.authorization.k8s.io/v1beta1", Kind: "RoleList"},
	}}
	testOpts := DiscoveryClient{
		ClientSet: clientset,
		Instance:  instance,
		DiscoveryClient: &preferredResourcesDiscovery{
			DiscoveryInterface: &discoveryFake.FakeDiscovery{},
			resources: []*metav1.APIResourceList{
				{GroupVersion: "apps/v1", APIResources: []metav1.APIResource{{Name: "deployments", Kind: "Deployment", Namespaced: true}}},
				{GroupVersion: "rbac.authorization.k8s.io/v1", APIResources: []metav1.APIResource{{Name: "roles", Kind: "Role", Namespaced: true}}},
				{GroupVersion: "v1", APIResources: []metav1.APIResource{{Name: "configmaps", Kind: "ConfigMap", Namespaced: true}}},
			},
		},
	}

	assert.NoError(t, testOpts.GetApiResources())

	listed := []string{}
	for _, action := range clientset.Actions() {
		assert.Equal(t, "list", action.GetVerb())
		listed = append(listed, action.GetResource().Resource)
	}
	assert.Equal(t, []string{"deployments", "roles"}, listed)

	if assert.Len(t, instance.Outputs, 1) {
		assert.Equal(t, "web", instance.Outputs[0].Name)
		assert.Equal(t, extensionsDeployment, *instance.Outputs[0].APIVersion)
	}
}
