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
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

// Version is an apiVersion and a flag for deprecation
type Version struct {
	// Name is the name of the api version
	Name string `json:"version" yaml:"version"`
	// Kind is the kind of object associated with this version
	Kind string `json:"kind" yaml:"kind"`
	// DeprecatedIn is a string that indicates what version the api is deprecated in
	// an empty string indicates that the version is not deprecated
	DeprecatedIn string `json:"deprecated-in" yaml:"deprecated-in"`
	// RemovedIn denotes the version that the api was actually removed in
	// An empty string indicates that the version has not been removed yet
	RemovedIn string `json:"removed-in" yaml:"removed-in"`
	// ReplacementAPI is the apiVersion that replaces the deprecated one
	ReplacementAPI string `json:"replacement-api" yaml:"replacement-api"`
	// Component is the component associated with this version
	Component string `json:"-" yaml:"-"`
}

func checkVersion(stub *Stub) *Version {
	for _, version := range VersionList {
		// We allow empty kinds to deprecate whole APIs.
		if version.Kind == "" || version.Kind == stub.Kind {
			if version.Name == stub.APIVersion {
				if version.Kind == "" {
					version.Kind = stub.Kind
				}

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
				output.Namespace = stub.Metadata.Namespace
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

// IsDeprecatedIn returns true if the version is deprecated in the applicable targetVersion
// Will return false if the targetVersion passed is not a valid semver string
func (v *Version) isDeprecatedIn(targetVersions map[string]string) bool {
	for component, targetVersion := range targetVersions {
		if !semver.IsValid(targetVersion) {
			klog.V(3).Infof("targetVersion %s for %s is not valid semVer", targetVersion, component)
			return false
		}
	}

	if v.DeprecatedIn == "" {
		return false
	}

	targetVersion, ok := targetVersions[v.Component]
	if !ok {
		klog.V(3).Infof("targetVersion missing for component %s", v.Component)
		return false
	}

	comparison := semver.Compare(targetVersion, v.DeprecatedIn)
	return comparison >= 0
}

// IsRemovedIn returns true if the version is deprecated in the applicable targetVersion
// Will return false if the targetVersion passed is not a valid semver string
func (v *Version) isRemovedIn(targetVersions map[string]string) bool {
	for component, targetVersion := range targetVersions {
		if !semver.IsValid(targetVersion) {
			klog.V(3).Infof("targetVersion %s for %s is not valid semVer", targetVersion, component)
			return false
		}
	}

	if v.RemovedIn == "" {
		return false
	}

	targetVersion, ok := targetVersions[v.Component]
	if !ok {
		klog.V(3).Infof("targetVersion missing for component %s", v.Component)
		return false
	}

	comparison := semver.Compare(targetVersion, v.RemovedIn)
	return comparison >= 0
}
