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

package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

// Stub is a stub of a Kubernetes manifest that has just the name and apiVersion
type Stub struct {
	Kind       string   `json:"kind" yaml:"kind"`
	APIVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Metadata   StubMeta `json:"metadata" yaml:"metadata"`
	Items      []Stub   `json:"items" yaml:"items"`
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
	// ReplacementAvailableIn is the version in which the replacement api is available
	ReplacementAvailableIn string `json:"replacement-available-in" yaml:"replacement-available-in"`
	// Component is the component associated with this version
	Component string `json:"component" yaml:"component"`
}

// VersionFile is a file with a list of deprecated versions
type VersionFile struct {
	DeprecatedVersions []Version         `json:"deprecated-versions" yaml:"deprecated-versions"`
	TargetVersions     map[string]string `json:"target-versions,omitempty" yaml:"target-versions,omitempty"`
}

func (instance *Instance) checkVersion(stub *Stub) *Version {
	for _, version := range instance.DeprecatedVersions {
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
func (instance *Instance) IsVersioned(data []byte) ([]*Output, error) {
	var outputs []*Output
	stubs, err := containsStub(data)
	if err != nil {
		return nil, err
	}
	if len(stubs) > 0 {
		for _, stub := range stubs {
			var output Output
			version := instance.checkVersion(stub)
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
	return nil, nil
}

// containsStub checks to see if a []byte has a stub in it
func containsStub(data []byte) ([]*Stub, error) {
	klog.V(10).Infof("\n%s", string(data))
	stub, err := jsonToStub(data)
	if err != nil {
		klog.V(8).Infof("invalid json: %s, trying yaml", err.Error())
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
	expandList(&stubs, stub)
	return stubs, nil
}

func yamlToStub(data []byte) ([]*Stub, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	var stubs []*Stub
	var tError *yaml.TypeError
	var errs []error
	for {
		stub := &Stub{}
		err := decoder.Decode(stub)
		if err != nil {
			if err == io.EOF {
				break
			}
			if errors.As(err, &tError) {
				klog.V(2).Infof("skipping for invalid yaml in manifest: %s", err)
				errs = append(errs, err)
				continue
			}
			return stubs, err
		}
		expandList(&stubs, stub)
	}
	if stubs == nil && len(errs) > 0 {
		return nil, fmt.Errorf("one or more errors parsing yaml resulted in no versions found: %v", errs)
	}
	return stubs, nil
}

// expandList checks if we have a List manifest.
// If it is the case, the manifests inside are expanded, otherwise we just return the single manifest
func expandList(stubs *[]*Stub, currentStub *Stub) {
	if len(currentStub.Items) > 0 {
		klog.V(5).Infof("found a list with %d items, attempting to expand", len(currentStub.Items))
		for _, stub := range currentStub.Items {
			currentItem := stub
			*stubs = append(*stubs, &currentItem)
		}
	} else {
		*stubs = append(*stubs, currentStub)
	}
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

// isReplacementAvailableIn returns true if the replacement api is available in the applicable targetVersion
// Will return false if the targetVersion passed is not a valid semver string
func (v *Version) isReplacementAvailableIn(targetVersions map[string]string) bool {
	for component, targetVersion := range targetVersions {
		if !semver.IsValid(targetVersion) {
			klog.V(3).Infof("targetVersion %s for %s is not valid semVer", targetVersion, component)
			return false
		}
	}

	if v.ReplacementAvailableIn == "" {
		return true
	}

	targetVersion, ok := targetVersions[v.Component]
	if !ok {
		klog.V(3).Infof("targetVersion missing for component %s", v.Component)
		return false
	}

	comparison := semver.Compare(targetVersion, v.ReplacementAvailableIn)
	return comparison >= 0
}

// PrintVersionList prints out the list of versions
// in a specific format
func (instance *Instance) PrintVersionList(outputFormat string) error {
	switch outputFormat {
	case "normal", "wide":
		err := instance.printVersionsTabular()
		if err != nil {
			return err
		}
	case "json":
		versionFile := VersionFile{
			DeprecatedVersions: instance.DeprecatedVersions,
			TargetVersions:     instance.TargetVersions,
		}
		data, err := json.Marshal(versionFile)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "yaml":
		versionFile := VersionFile{
			DeprecatedVersions: instance.DeprecatedVersions,
			TargetVersions:     instance.TargetVersions,
		}
		data, err := yaml.Marshal(versionFile)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	default:
		errText := "The output format must one of (normal|wide|json|yaml)"
		fmt.Println(errText)
		return fmt.Errorf(errText)
	}
	return nil
}

func (instance *Instance) printVersionsTabular() error {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 15, 2, padChar, 0)

	if !instance.NoHeaders {
		fmt.Fprintln(w, "KIND\t NAME\t DEPRECATED IN\t REMOVED IN\t REPLACEMENT\t REPL AVAIL IN\t COMPONENT\t")
	}

	for _, version := range instance.DeprecatedVersions {
		deprecatedIn := version.DeprecatedIn
		if deprecatedIn == "" {
			deprecatedIn = "n/a"
		}
		removedIn := version.RemovedIn
		if removedIn == "" {
			removedIn = "n/a"
		}

		replacementAPI := version.ReplacementAPI
		if replacementAPI == "" {
			replacementAPI = "n/a"
		}

		replacementAvailableIn := version.ReplacementAvailableIn
		if replacementAvailableIn == "" {
			replacementAvailableIn = "n/a"
		}

		_, _ = fmt.Fprintf(w, "%s\t %s\t %s\t %s\t %s\t %s\t %s\t\n", version.Kind, version.Name, deprecatedIn, removedIn, replacementAPI, replacementAvailableIn, version.Component)
	}
	err := w.Flush()
	if err != nil {
		return err
	}
	return nil
}

// UnMarshalVersions reads data from a versions file and returns the versions
// If included, it will also return the map of targetVersions
func UnMarshalVersions(data []byte) ([]Version, map[string]string, error) {
	versionFile := &VersionFile{}
	err := yaml.Unmarshal(data, versionFile)
	if err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal versions file from data: %s", err.Error())
	}
	return versionFile.DeprecatedVersions, versionFile.TargetVersions, nil

}

// GetDefaultVersionList gets the default versions from the versions.yaml file
func GetDefaultVersionList(versionFileData []byte) ([]Version, map[string]string, error) {
	defaultVersions, defaultTargetVersions, err := UnMarshalVersions(versionFileData)
	if err != nil {
		return nil, nil, err
	}
	return defaultVersions, defaultTargetVersions, nil
}

// CombineAdditionalVersions adds additional versions into the defaults. If the additional versions
// contain any that already exist in the defaults, return an error
func CombineAdditionalVersions(additional []Version, defaults []Version) ([]Version, error) {
	returnList := defaults
	for _, version := range additional {
		klog.V(3).Infof("attempting to combine into defaults: %v", version)
		if version.isContainedIn(defaults) {
			return nil, fmt.Errorf("duplicate cannot be added to defaults: %s %s", version.Kind, version.Name)
		}
		returnList = append(returnList, version)
	}
	return returnList, nil
}

func (v Version) isContainedIn(versionList []Version) bool {
	for _, version := range versionList {
		if isDuplicate(v, version) {
			return true
		}
	}
	return false
}

func isDuplicate(a Version, b Version) bool {
	if a.Kind == b.Kind {
		if a.Name == b.Name {
			return true
		}
	}
	return false
}
