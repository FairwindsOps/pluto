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
	"encoding/json"
	"fmt"

	helmstoragev2 "helm.sh/helm/pkg/storage"
	driverv2 "helm.sh/helm/pkg/storage/driver"
	helmstoragev3 "helm.sh/helm/v3/pkg/storage"
	driverv3 "helm.sh/helm/v3/pkg/storage/driver"

	"github.com/fairwindsops/pluto/pkg/api"
)

// Helm represents all current releases that we can find in the cluster
type Helm struct {
	Releases []*Release
	Outputs  []*api.Output
	Version  string
	Kube     *kube
}

type Release struct {
	Name     string `json:"name"`
	Chart    *Chart `json:"chart"`
	Manifest string `json:"manifest"`
}

type Chart struct {
	Metadata *ChartMeta `json:"metadata"`
}

type ChartMeta struct {
	Name    string `json:"name"`
	Version string `json:"version"`
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
		err = h.getReleasesVersionTwo()
	case "3":
		err = h.getReleasesVersionThree()
	default:
		err = fmt.Errorf("helm version either not specified or incorrect (use 2 or 3)")
	}
	return err
}

// getReleasesVersionTwo retrieves helm 2 releases from ConfigMaps
func (h *Helm) getReleasesVersionTwo() error {
	if h.Version != "2" {
		return fmt.Errorf("helm 2 function called without helm 2 version set")
	}
	hcm := driverv2.NewConfigMaps(h.Kube.Client.CoreV1().ConfigMaps(""))
	helmClient := helmstoragev2.Init(hcm)
	list, err := helmClient.ListDeployed()
	if err != nil {
		return err
	}
	for _, release := range list {
		deployed, err := helmClient.Deployed(release.Name)
		if err != nil {
			return fmt.Errorf("error determining most recent deployed for '%s'\n   %w", release.Name, err)
		}
		if release.Version != deployed.Version {
			continue
		}
		jsonRel, err := json.Marshal(release)
		if err != nil {
			return fmt.Errorf("error marshaling release '%s'\n   %w", release.Name, err)
		}
		rel, err := marshalToRelease(jsonRel)
		if err != nil {
			return fmt.Errorf("error unmarshaling release '%s'\n   %w", release.Name, err)
		}
		h.Releases = append(h.Releases, rel)
	}
	if err := h.findVersions(); err != nil {
		return err
	}
	return nil
}

// getReleasesVersionThree retrieves helm 3 releases from Secrets
func (h *Helm) getReleasesVersionThree() error {
	if h.Version != "3" {
		return fmt.Errorf("helm 3 function called without helm 3 version set")
	}
	hs := driverv3.NewSecrets(h.Kube.Client.CoreV1().Secrets(""))
	helmClient := helmstoragev3.Init(hs)
	list, err := helmClient.ListDeployed()
	if err != nil {
		return err
	}
	for _, release := range list {
		deployed, err := helmClient.Deployed(release.Name)
		if err != nil {
			return fmt.Errorf("error determining most recent deployed for '%s'\n   %w", release.Name, err)
		}
		if release.Version != deployed.Version {
			continue
		}
		jsonRel, err := json.Marshal(release)
		if err != nil {
			return fmt.Errorf("error marshaling release '%s'\n   %w", release.Name, err)
		}
		rel, err := marshalToRelease(jsonRel)
		if err != nil {
			return fmt.Errorf("error unmarshaling release '%s'\n   %w", release.Name, err)
		}
		h.Releases = append(h.Releases, rel)
	}
	if err := h.findVersions(); err != nil {
		return err
	}
	return nil
}

func (h *Helm) findVersions() error {
	fmt.Println("done")
	for _, release := range h.Releases {
		outList, err := checkForAPIVersion([]byte(release.Manifest))
		if err != nil {
			return fmt.Errorf("error parsing release '%s'\n   %w", release.Name, err)
		}
		for _, out := range outList {
			out.Name = release.Name + "/" + out.Name
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

// marshalToRelease casts a marshals release data into our Release type so we have a common type regardless of helm version
func marshalToRelease(jsonRel []byte) (*Release, error) {
	var ret = new(Release)
	err := json.Unmarshal(jsonRel, ret)
	return ret, err
}
