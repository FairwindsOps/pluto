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

	helmstoragev2 "helm.sh/helm/pkg/storage"
	driverv2 "helm.sh/helm/pkg/storage/driver"
	helmstoragev3 "helm.sh/helm/v3/pkg/storage"
	driverv3 "helm.sh/helm/v3/pkg/storage/driver"

	"github.com/fairwindsops/pluto/pkg/api"
)

// Helm represents all current releases that we can find in the cluster
type Helm struct {
	CurrentReleases []*Release
	Outputs         []*api.Output
	Version         string
	Kube            *kube
}

// Release is a single version of a chart release
type Release struct {
	Name     string `json:"metadata.name" yaml:"metadata.name"`
	Manifest string `json:"manifest,omitempty"`
}

// NewHelm returns a basic helm struct with the version of helm requested
func NewHelm(version string) *Helm {
	return &Helm{
		Version: version,
		Kube:    getConfigInstance(),
	}
}

// FindVersions is the primary method in the package. It ties together all the functionality
func (h *Helm) FindVersions() error {
	var err error

	switch h.Version {
	case "2":
		err = h.getManifestsVersionTwo()
	case "3":
		err = h.getManifestsVersionThree()
	default:
		err = fmt.Errorf("helm version either not specified or incorrect (use 2 or 3)")
	}

	return err
}

// getManifestsVersionTwo retrieves helm 2 manifests from ConfigMaps
func (h *Helm) getManifestsVersionTwo() error {
	if h.Version != "2" {
		return fmt.Errorf("helm 2 function called without helm 2 version set")
	}
	hcm := driverv2.NewConfigMaps(h.Kube.Client.CoreV1().ConfigMaps(""))
	helmClient := helmstoragev2.Init(hcm)
	list, _ := helmClient.ListDeployed()
	for _, release := range list {
		outList, err := checkForAPIVersion([]byte(release.Manifest))
		if err != nil {
			return err
		}
		h.Outputs = append(h.Outputs, outList...)
	}
	return nil
}

// getManifestsVersionThree retrieves helm 3 manifests from Secrets
func (h *Helm) getManifestsVersionThree() error {
	if h.Version != "3" {
		return fmt.Errorf("helm 3 function called without helm 3 version set")
	}
	hs := driverv3.NewSecrets(h.Kube.Client.CoreV1().Secrets(""))
	helmClient := helmstoragev3.Init(hs)
	list, _ := helmClient.ListDeployed()
	for _, release := range list {
		outList, err := checkForAPIVersion([]byte(release.Manifest))
		if err != nil {
			return err
		}
		h.Outputs = append(h.Outputs, outList...)
	}
	return nil
}

// checkForAPIVersion calls the api pkg to parse our releases for deprecated APIs
func checkForAPIVersion(manifest []byte) ([]*api.Output, error) {
	outputs, err := api.IsVersioned(manifest)
	if err != nil {
		return nil, err
	}
	if len(outputs) < 1 {
		return nil, nil
	}
	return outputs, nil
}
