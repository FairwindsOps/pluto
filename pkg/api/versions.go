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
	// Deprecated is a boolean that indicates if the apiVersion is deprecated
	Deprecated bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
}

// VersionList is a set of apiVersions and if they are deprecated or not.
// TODO: Find a better way to generate this.
// Currently using the list for 1.16 from here: https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/
var VersionList = []Version{
	// Deployments
	{"apps/v1", "Deployment", false},
	{"extensions/v1beta1", "Deployment", true},
	{"apps/v1beta2", "Deployment", true},
	{"apps/v1beta1", "Deployment", true},

	// StatefulSet
	{"apps/v1beta1", "StatefulSet", true},
	{"apps/v1beta2", "StatefulSet", true},

	// NetworkPolicy
	{"networking.k8s.io/v1", "NetworkPolicy", false},
	{"extensions/v1beta1", "NetworkPolicy", true},

	// DaemonSet
	{"apps/v1beta2", "DaemonSet", true},
	{"extensions/v1beta1", "DaemonSet", true},

	// PodSecurityPolicy
	{"policy/v1beta1", "PodSecurityPolicy", false},
	{"extensions/v1beta1", "PodSecurityPolicy", true},

	// ReplicaSet
	{"extensions/v1beta1", "ReplicaSet", true},
	{"apps/v1beta1", "ReplicaSet", true},
	{"apps/v1beta2", "ReplicaSet", true},
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
		klog.V(8).Infof("not json: %s", err.Error())
	} else {
		return stub, nil
	}
	stub, err = yamlToStub(data)
	if err != nil {
		klog.V(8).Infof("not yaml: %s", err.Error())
	} else {
		return stub, nil
	}
	return nil, fmt.Errorf("no matches for stub, error: %v", err)
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
