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
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	"github.com/fairwindsops/api-version-finder/pkg/api"
)

type Helm struct {
	CurrentReleases []*Release
	Outputs         []*api.Output
	Version         string
}

type Release struct {
	Name     string `json:"metadata.name" yaml:"metadata.name"`
	Manifest string `json:"manifest,omitempty"`
}

func NewHelm(version string) *Helm {
	return &Helm{
		Version: version,
	}
}

func (h *Helm) FindVersions() error {
	var err error
	var k = getConfigInstance()
	switch h.Version {
	case "2":
		err = h.getManifestsVersionTwo(k)
	case "3":
		err = h.getManifestsVersionThree(k)
	default:
		err = fmt.Errorf("helm version either not specified or incorrect (use 2 or 3)")
	}

	return err
}

// getManifestsVersionTwo retrieves helm 2 manifests from ConfigMaps
func (h *Helm) getManifestsVersionTwo(k *Kube) error {
	return fmt.Errorf("helm 2 check not implemented")
}

// getManifestsVersionThree retrieves helm 3 manifests from Secrets
func (h *Helm) getManifestsVersionThree(k *Kube) error {
	secrets, err := k.Client.CoreV1().Secrets("").List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	var releases []*Release
	for _, secret := range secrets.Items {
		if secret.Type != "helm.sh/release.v1" {
			continue
		}
		thisRelease, err := decodeReleaseSecret(string(secret.Data["release"]))
		if err != nil {
			klog.Error("Failed parsing secret release data for helm 3: ", err)
			return err
		}
		thisRelease.Name = secret.Name
		releases = append(releases, thisRelease)
	}
	err = h.setCurrentReleases(releases)
	if err != nil {
		return err
	}
	for _, release := range h.CurrentReleases {
		outList, err := checkForAPIVersion([]byte(release.Manifest))
		if err != nil {
			return err
		}
		h.Outputs = append(h.Outputs, outList...)
	}
	return nil
}

// setCurrentReleases parses all the releases we found and makes sure we only return the most recent
// which should be the one deployed currently
func (h *Helm) setCurrentReleases(allReleases []*Release) error {
	found := map[string]*Release{}
	regex := regexp.MustCompile(`^(sh\.helm\.release)\.(v[0-9]+)\.([\w\-]+)\.v([0-9]+)$`)
	for _, release := range allReleases {
		parts := regex.FindStringSubmatch(release.Name)
		releaseName := parts[3]
		thisVersion, err := strconv.Atoi(parts[4])
		if err != nil {
			return err
		}
		if v, ok := found[releaseName]; ok {
			storedParts := regex.FindStringSubmatch(v.Name)
			storedVersion, err := strconv.Atoi(storedParts[4])
			if err != nil {
				return err
			}
			if thisVersion > storedVersion {
				found[releaseName] = release
			}
		} else {
			found[releaseName] = release
		}
	}
	for _, v := range found {
		h.CurrentReleases = append(h.CurrentReleases, v)
	}
	return nil
}

// This function is ripped straight out of the helm 3 codebase with slight modification
// https://github.com/helm/helm/blob/193850a9e2c509acf1a499d98e8d23c12c134f11/pkg/storage/driver/util.go#L56-L84
func decodeReleaseSecret(data string) (*Release, error) {
	var magicGzip = []byte{0x1f, 0x8b, 0x08}
	var b64 = base64.StdEncoding
	// base64 decode string
	b, err := b64.DecodeString(data)
	if err != nil {
		return nil, err
	}

	// For backwards compatibility with releases that were stored before
	// compression was introduced we skip decompression if the
	// gzip magic header is not found
	if bytes.Equal(b[0:3], magicGzip) {
		r, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		b2, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		b = b2
	}

	var secOut Release
	// unmarshal release object bytes
	if err := json.Unmarshal(b, &secOut); err != nil {
		return nil, err
	}
	return &secOut, nil
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
