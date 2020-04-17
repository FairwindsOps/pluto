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

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v3"
	"k8s.io/klog"
)

// Stub is a stub of a Kubernetes manifest that has just the name and apiVersion
type Stub struct {
	Kind       string   `json:"kind" yaml:"kind"`
	APIVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Metadata   StubMeta `json:"metadata" yaml:"metadata"`
}

// StubMeta will catch kube resource metadata
type StubMeta struct {
	Name string `json:"name" yaml:"name"`
}

// Output is a thing that has an apiVersion in it
type Output struct {
	Name       string   `json:"file,omitempty" yaml:"file,omitempty"`
	APIVersion *Version `json:"api,omitempty" yaml:"api,omitempty"`
}

// Version is an apiVersion and a flag for deprecation
type Version struct {
	// Name is the name of the api version
	Name string `json:"version,omitempty" yaml:"version,omitempty"`
	// Kind is the kind of object associated with this version
	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`
	// DeprecatedIn is a string that indicates what version the api is deprecated in
	// an empty string indicates that the version is not deprecated
	DeprecatedIn string `json:"deprecated-in,omitempty" yaml:"deprecated-in,omitempty"`
	// RemovedIn denotes the version that the api was actually removed in
	// An empty string indicates that the version has not been removed yet
	RemovedIn string `json:"removed-in,omitempty" yaml:"removed-in,omitempty"`
	// ReplacementAPI is the apiVersion that replaces the deprecated one
	ReplacementAPI string `json:"replacement-api,omitempty" yaml:"replacement-api,omitempty"`
}

// VersionList is a set of apiVersions and if they are deprecated or not.
// TODO: Find a better way to generate this.
// Currently using the list for 1.16 from here: https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/
var VersionList = []Version{
	// Deployments
	{"apps/v1", "Deployment", "", "", ""},
	{"extensions/v1beta1", "Deployment", "v1.9.0", "v1.16.0", "apps/v1"},
	{"apps/v1beta2", "Deployment", "v1.9.0", "v1.16.0", "apps/v1"},
	{"apps/v1beta1", "Deployment", "v1.9.0", "v1.16.0", "apps/v1"},

	// StatefulSet
	{"apps/v1beta1", "StatefulSet", "v1.9.0", "v1.16.0", "apps/v1"},
	{"apps/v1beta2", "StatefulSet", "v1.9.0", "v1.16.0", "apps/v1"},

	// NetworkPolicy
	{"networking.k8s.io/v1", "NetworkPolicy", "", "", ""},
	{"extensions/v1beta1", "NetworkPolicy", "v1.9.0", "v1.16.0", "networking.k8s.io/v1"},

	// DaemonSet
	{"apps/v1beta2", "DaemonSet", "v1.9.0", "v1.16.0", "apps/v1"},
	{"extensions/v1beta1", "DaemonSet", "v1.9.0", "v1.16.0", "apps/v1"},

	// PodSecurityPolicy
	{"policy/v1beta1", "PodSecurityPolicy", "", "", ""},
	{"extensions/v1beta1", "PodSecurityPolicy", "v1.10.0", "v1.16.0", "policy/v1beta1"},

	// ReplicaSet
	{"extensions/v1beta1", "ReplicaSet", "", "v1.16.0", "apps/v1"},
	{"apps/v1beta1", "ReplicaSet", "", "v1.16.0", "apps/v1"},
	{"apps/v1beta2", "ReplicaSet", "", "v1.16.0", "apps/v1"},

	// MutatingWebhookConfiguration
	{"admissionregistration.k8s.io/v1beta1", "MutatingWebhookConfiguration", "v1.16.0", "v1.19.0", "admissionregistration.k8s.io/v1"},

	// CustomResourceDefinition
	{"apiextensions.k8s.io/v1beta1", "CustomResourceDefinition", "v1.19.0"},

	// PriorityClass
	{"scheduling.k8s.io/v1beta1", "PriorityClass", "v1.17.0"},
	{"scheduling.k8s.io/v1alpha1", "PriorityClass", "v1.17.0"},
	{"apiextensions.k8s.io/v1beta1", "CustomResourceDefinition", "v1.16.0", "v1.19.0", "apiextensions.k8s.io/v1"},

	// Rbac
	{"rbac.authorization.k8s.io/v1alpha1", "ClusterRoleBinding", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},
	{"rbac.authorization.k8s.io/v1alpha1", "ClusterRole", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},
	{"rbac.authorization.k8s.io/v1alpha1", "ClusterRoleBindingList", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},
	{"rbac.authorization.k8s.io/v1alpha1", "ClusterRoleList", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},
	{"rbac.authorization.k8s.io/v1alpha1", "Role", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},
	{"rbac.authorization.k8s.io/v1alpha1", "RoleBinding", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},
	{"rbac.authorization.k8s.io/v1alpha1", "RoleList", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},
	{"rbac.authorization.k8s.io/v1alpha1", "RoleBindingList", "v1.17.0", "v1.20.0", "rbac.authorization.k8s.io/v1"},

	//CSINode
	{"storage.k8s.io/v1beta1", "CSINode", "v1.17.0", "", ""},
>>>>>>> ac5df95... Fix #46 add remaining v1.17 deprecations
}

func checkVersion(stub *Stub) *Version {
	for _, version := range VersionList {
		if version.Kind == stub.Kind {
			if version.Name == stub.APIVersion {
				return &version
			}
		}
	}
	return nil
}

// IsVersioned returns a version if the file data sent
// can be unmarshaled into a stub and matches a known
// version in the VersionList
func IsVersioned(data []byte) ([]*Output, error) {
	var outputs []*Output
	stubs, err := containsStub(data)
	if err != nil {
		return nil, err
	}
	if len(stubs) > 0 {
		for _, stub := range stubs {
			var output Output
			version := checkVersion(stub)
			if version != nil {
				output.Name = stub.Metadata.Name
				output.APIVersion = version
			} else {
				continue
			}
			outputs = append(outputs, &output)
		}
		return outputs, nil
	}
	return nil, fmt.Errorf("no version found in data")
}

// containsStub checks to see if a []byte has a stub in it
func containsStub(data []byte) ([]*Stub, error) {
	klog.V(10).Infof("\n%s", string(data))
	stub, err := jsonToStub(data)
	if err != nil {
		klog.V(8).Infof("invalid json: %s", err.Error())
	} else {
		return stub, nil
	}
	stub, err = yamlToStub(data)
	if err != nil {
		klog.V(8).Infof("invalid yaml: %s", err.Error())
	} else {
		return stub, nil
	}
	return nil, err
}

func jsonToStub(data []byte) ([]*Stub, error) {
	var stubs []*Stub
	stub := &Stub{}
	err := json.Unmarshal(data, stub)
	if err != nil {
		return nil, err
	}
	stubs = append(stubs, stub)
	return stubs, nil
}

func yamlToStub(data []byte) ([]*Stub, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	var stubs []*Stub
	for {
		stub := &Stub{}
		err := decoder.Decode(stub)
		if err != nil {
			if err == io.EOF {
				break
			}
			return stubs, err
		}
		stubs = append(stubs, stub)
	}
	return stubs, nil
}

// IsDeprecatedIn returns true if the version is deprecated in the targetVersion
// Will return false if the targetVersion passed is not a valid semver string
func (v *Version) IsDeprecatedIn(targetVersion string) bool {
	if !semver.IsValid(targetVersion) {
		klog.V(3).Infof("targetVersion %s is not valid semVer", targetVersion)
		return false
	}
	if v.DeprecatedIn == "" {
		return false
	}
	comparison := semver.Compare(targetVersion, v.DeprecatedIn)
	return comparison >= 0
}

// IsRemovedIn returns true if the version is deprecated in the targetVersion
// Will return false if the targetVersion passed is not a valid semver string
func (v *Version) IsRemovedIn(targetVersion string) bool {
	if !semver.IsValid(targetVersion) {
		klog.V(3).Infof("targetVersion %s is not valid semVer", targetVersion)
		return false
	}
	if v.RemovedIn == "" {
		return false
	}
	comparison := semver.Compare(targetVersion, v.RemovedIn)
	return comparison >= 0
}
