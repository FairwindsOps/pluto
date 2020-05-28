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
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/fairwindsops/pluto/pkg/api"
)

var (
	helmConfigMap = v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "helmtest.v1",
			Labels: map[string]string{
				"MODIFIED_AT": "1585848156",
				"NAME":        "helmtest",
				"OWNER":       "TILLER",
				"STATUS":      "DEPLOYED",
				"VERSION":     "1",
			},
		},
		Data: map[string]string{
			"release": "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",
		},
	}
	helmConfigMap2 = v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "helmtest2.v2",
			Labels: map[string]string{
				"MODIFIED_AT": "1589922959",
				"NAME":        "kindled-toucan",
				"OWNER":       "TILLER",
				"STATUS":      "DEPLOYED",
				"VERSION":     "1",
			},
		},
		Data: map[string]string{
			"release": "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",
		},
	}
	helmSecret = v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sh.helm.release.v1.helmtest.v1",
			Labels: map[string]string{
				"createdAt": "1585859667",
				"name":      "helmtest",
				"owner":     "helm",
				"status":    "deployed",
				"version":   "1",
			},
		},
		Data: map[string][]byte{
			"release": []byte("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"),
		},
	}
	helmSecret2 = v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sh.helm.release.v1.helmtest.v1",
			Labels: map[string]string{
				"MODIFIED_AT": "1585859667",
				"NAME":        "helmtest",
				"OWNER":       "TILLER",
				"STATUS":      "DEPLOYED",
				"VERSION":     "1",
			},
		},
		Data: map[string][]byte{
			"release": []byte("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"),
		},
	}
	wantOutput = []*api.Output{
		{
			Name:      "helmtest/helmtest-helmchartest-v1beta1",
			Namespace: "default",
			APIVersion: &api.Version{
				Name:           "extensions/v1beta1",
				Kind:           "Deployment",
				DeprecatedIn:   "v1.9.0",
				RemovedIn:      "v1.16.0",
				ReplacementAPI: "apps/v1",
			},
		},
		{
			Name:      "helmtest/helmtest-helmchartest",
			Namespace: "default",
			APIVersion: &api.Version{
				Name:           "apps/v1",
				Kind:           "Deployment",
				DeprecatedIn:   "",
				RemovedIn:      "",
				ReplacementAPI: "",
			},
		},
	}
	wantOutputNamespaced = []*api.Output{
		{
			Name:      "kindled-toucan/kindled-toucan-basic-demo",
			Namespace: "demo1",
			APIVersion: &api.Version{
				Name:           "apps/v1",
				Kind:           "Deployment",
				DeprecatedIn:   "",
				RemovedIn:      "",
				ReplacementAPI: "",
			},
		},
	}
)

func newMockHelm(version, store, namespace string) *Helm {
	return &Helm{
		Version:   version,
		Namespace: namespace,
		Kube:      getMockConfigInstance(),
		Store:     store,
	}
}

func newBadKubeClient() (k *kube) {
	conf := new(rest.Config)
	conf.Host = "127.0.0.1:9999"
	k = new(kube)
	k.Client, _ = kubernetes.NewForConfig(conf)
	return
}

func Test_checkForAPIVersion(t *testing.T) {
	tests := []struct {
		name     string
		manifest []byte
		want     []*api.Output
		wantErr  bool
	}{
		{
			name:     "empty",
			manifest: []byte{},
			want:     []*api.Output{{}},
			wantErr:  true,
		},
		{
			name:     "got version",
			manifest: []byte("apiVersion: apps/v1\nkind: Deployment"),
			want:     []*api.Output{{APIVersion: &api.Version{Name: "apps/v1", Kind: "Deployment", DeprecatedIn: "", RemovedIn: "", ReplacementAPI: ""}}},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkForAPIVersion(tt.manifest)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, tt.want, got)

		})
	}
}

func TestHelm_getManifestsVersionTwo(t *testing.T) {
	tests := []struct {
		name        string
		helmVersion string
		store       string
		namespace   string
		wantErr     bool
		errMessage  string
		configMap   *v1.ConfigMap
		secret      *v1.Secret
		want        []*api.Output
	}{
		{
			name:        "three - error",
			helmVersion: "3",
			wantErr:     true,
			errMessage:  "helm 2 function called without helm 2 version set",
		},
		{
			name:        "helm 2 valid configmap",
			helmVersion: "2",
			store:       "configmaps",
			configMap:   &helmConfigMap,
			want:        wantOutput,
		},
		{
			name:        "helm 2 valid secret",
			helmVersion: "2",
			store:       "secrets",
			secret:      &helmSecret2,
			want:        wantOutput,
		},
		{
			name:        "helm 2 other namespace",
			namespace:   "demo1",
			helmVersion: "2",
			store:       "configmaps",
			configMap:   &helmConfigMap,
			want:        nil,
		},
		{
			name:        "helm 2 target namespace",
			namespace:   "demo1",
			helmVersion: "2",
			store:       "configmaps",
			configMap:   &helmConfigMap2,
			want:        wantOutputNamespaced,
		},
	}
	for _, tt := range tests {
		h := newMockHelm(tt.helmVersion, tt.store, tt.namespace)
		if tt.configMap != nil {
			_, err := h.Kube.Client.CoreV1().ConfigMaps(tt.namespace).Create(tt.configMap)
			if err != nil {
				t.Errorf("failed putting configMap in mocked kube. test: %s", tt.name)
			}
		}
		if tt.secret != nil {
			_, err := h.Kube.Client.CoreV1().Secrets(tt.namespace).Create(tt.secret)
			if err != nil {
				t.Errorf("failed putting secret in mocked kube. test: %s", tt.name)
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			err := h.getReleasesVersionTwo()
			if tt.wantErr {
				assert.EqualError(t, err, tt.errMessage)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, h.Outputs)
		})
	}
}

func TestHelm_getManifestsVersionThree(t *testing.T) {
	tests := []struct {
		name        string
		helmVersion string
		wantErr     bool
		errMessage  string
		secret      *v1.Secret
		want        []*api.Output
	}{
		{
			name:        "two - error",
			helmVersion: "2",
			wantErr:     true,
			errMessage:  "helm 3 function called without helm 3 version set",
		},
		{
			name:        "helm 3 valid",
			helmVersion: "3",
			secret:      &helmSecret,
			want:        wantOutput,
		},
	}

	for _, tt := range tests {
		h := newMockHelm(tt.helmVersion, "secrets", "")
		if tt.secret != nil {
			_, err := h.Kube.Client.CoreV1().Secrets("default").Create(tt.secret)
			if err != nil {
				t.Errorf("failed putting secret in mocked kube. test: %s", tt.name)
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			err := h.getReleasesVersionThree()
			if tt.wantErr {
				assert.EqualError(t, err, tt.errMessage)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, h.Outputs)
		})
	}
}

func TestHelm_getManifest_badClient(t *testing.T) {
	tests := []struct {
		name        string
		helmVersion string
		wantErr     bool
		errMessage  string
		secret      *v1.Secret
		want        []*api.Output
	}{
		{
			name:        "two - bad client",
			helmVersion: "2",
			wantErr:     true,
			errMessage:  "helm 3 function called without helm 3 version set",
		},
		{
			name:        "three - bad client",
			helmVersion: "3",
			wantErr:     true,
			errMessage:  "helm 3 function called without helm 3 version set",
		},
	}

	for _, tt := range tests {
		h := &Helm{
			Version: tt.helmVersion,
			Kube:    newBadKubeClient(),
		}
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch tt.helmVersion {
			case "2":
				err = h.getReleasesVersionTwo()
			case "3":
				err = h.getReleasesVersionThree()
			}
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "connect: connection refused")
				return
			}
		})
	}
}

func TestHelm_FindVersions(t *testing.T) {
	tests := []struct {
		name        string
		helmVersion string
		wantErr     bool
		errMessage  string
	}{
		// Only adding this one test case since the others generally cross into other functions.
		{"one - err", "1", true, "helm version either not specified or incorrect (use 2 or 3)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newMockHelm(tt.helmVersion, "", "secrets")
			err := h.FindVersions()
			if tt.wantErr {
				assert.EqualError(t, err, tt.errMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
