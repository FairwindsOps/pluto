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
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/releaseutil"
	helmstoragev3 "helm.sh/helm/v3/pkg/storage"
	driverv3 "helm.sh/helm/v3/pkg/storage/driver"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"

	"github.com/fairwindsops/pluto/v3/pkg/api"
)

// Helm represents all current releases that we can find in the cluster
type Helm struct {
	Releases  []*Release
	Version   string
	Kube      *kube
	Namespace string
	Store     string
	Instance  *api.Instance
}

// Release represents a single helm release
type Release struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Chart     *Chart `json:"chart"`
	Manifest  string `json:"manifest"`
}

// Chart represents a single helm chart
type Chart struct {
	Metadata *ChartMeta `json:"metadata"`
}

// ChartMeta is the metadata of a Helm chart
type ChartMeta struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// NewHelm returns a basic helm struct with the version of helm requested
func NewHelm(version, store, namespace string, instance *api.Instance) *Helm {
	return &Helm{
		Version:   version,
		Kube:      getConfigInstance(),
		Namespace: namespace,
		Store:     store,
		Instance:  instance,
	}
}

// NewHelmWithKubeClient returns a helm struct with version of helm requested
// and uses the passed in kube client as the cluster to operate on
func NewHelmWithKubeClient(version, store, namespace string, instance *api.Instance, kubeClient kubernetes.Interface) *Helm {
	return &Helm{
		Version: version,
		Kube: &kube{
			Client: kubeClient,
		},
		Namespace: namespace,
		Store:     store,
		Instance:  instance,
	}
}

// FindVersions is the primary method in the package. It ties together all the functionality
func (h *Helm) FindVersions() error {
	var err error

	switch h.Version {
	case "2":
		fmt.Println("DEPRECATION WARNING: Helm 2 will be deprecated very soon. Please migrate.")
		fmt.Println("please see https://helm.sh/blog/helm-v2-deprecation-timeline/ for more info")
		err = h.getReleasesVersionTwo()
	case "3":
		err = h.getReleasesVersionThree()
	default:
		err = fmt.Errorf("helm version either not specified or incorrect (use 2 or 3)")
	}
	return err
}

// getReleasesVersionTwo retrieves helm 2 releases from ConfigMaps or Secrets
func (h *Helm) getReleasesVersionTwo() error {
	var helmClient *helmstoragev2.Storage
	if h.Version != "2" {
		return fmt.Errorf("helm 2 function called without helm 2 version set")
	}
	switch h.Store {
	case "secrets":
		hs := driverv2.NewSecrets(h.Kube.Client.CoreV1().Secrets(h.Namespace))
		helmClient = helmstoragev2.Init(hs)
	case "configmaps":
		hcm := driverv2.NewConfigMaps(h.Kube.Client.CoreV1().ConfigMaps(""))
		helmClient = helmstoragev2.Init(hcm)
	default:
		return fmt.Errorf("helm-store should be configmap or secrets")
	}
	releases, err := helmClient.ListReleases()
	if err != nil {
		return err
	}
	for _, r := range releases {
		if h.Namespace != "" && r.Namespace != h.Namespace {
			continue
		}
		deployed, err := helmClient.Deployed(r.Name)
		if err != nil {
			klog.Infof("cannot determine most recent deployed for %s/%s - %s", r.Namespace, r.Name, err)
			continue
		}
		if r.Version != deployed.Version {
			continue
		}
		rel, err := helmToRelease(r)
		if err != nil {
			return fmt.Errorf("error converting helm r '%s/%s' to internal object\n   %w", r.Namespace, r.Name, err)
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
	hs := driverv3.NewSecrets(h.Kube.Client.CoreV1().Secrets(h.Namespace))
	helmClient := helmstoragev3.Init(hs)
	namespaces, err := h.Kube.Client.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	releases, err := helmClient.ListDeployed()
	if err != nil {
		return err
	}
	for _, namespace := range namespaces.Items {
		ns := namespace.Name
		if h.Namespace != "" && ns != h.Namespace {
			continue
		}
		filteredReleases := h.deployedReleasesPerNamespace(ns, releases)
		for _, r := range filteredReleases {
			rel, err := helmToRelease(r)
			if err != nil {
				return fmt.Errorf("error converting helm r '%s/%s' to internal object\n   %w", r.Namespace, r.Name, err)
			}
			h.Releases = append(h.Releases, rel)
		}
	}
	if err := h.findVersions(); err != nil {
		return err
	}
	return nil
}

func (h *Helm) deployedReleasesPerNamespace(namespace string, releases []*release.Release) []*release.Release {
	return releaseutil.All(deployed, relNamespace(namespace)).Filter(releases)
}

func deployed(rls *release.Release) bool {
	return rls.Info.Status == release.StatusDeployed
}

func relNamespace(ns string) releaseutil.FilterFunc {
	return func(rls *release.Release) bool {
		return rls.Namespace == ns
	}
}

func (h *Helm) findVersions() error {
	for _, r := range h.Releases {
		klog.V(2).Infof("parsing r %s", r.Name)
		outList, err := h.checkForAPIVersion([]byte(r.Manifest))
		if err != nil {
			return fmt.Errorf("error parsing r '%s/%s'\n   %w", r.Namespace, r.Name, err)
		}
		for _, out := range outList {
			out.Name = r.Name + "/" + out.Name
			out.Namespace = r.Namespace
		}
		h.Instance.Outputs = append(h.Instance.Outputs, outList...)

	}
	return nil
}

// checkForAPIVersion calls the api pkg to parse our releases for deprecated APIs
func (h *Helm) checkForAPIVersion(manifest []byte) ([]*api.Output, error) {
	outputs, err := h.Instance.IsVersioned(manifest)
	if err != nil {
		return nil, err
	}
	if len(outputs) < 1 {
		return nil, nil
	}
	return outputs, nil
}

func helmToRelease(helmRelease interface{}) (*Release, error) {
	jsonRel, err := json.Marshal(helmRelease)
	if err != nil {
		return nil, fmt.Errorf("error marshaling release: %s", err.Error())
	}
	return marshalToRelease(jsonRel)
}

// marshalToRelease marshals release data into the Pluto Release type so we have a common type regardless of helm version
func marshalToRelease(jsonRel []byte) (*Release, error) {
	var ret = new(Release)
	err := json.Unmarshal(jsonRel, ret)
	return ret, err
}
