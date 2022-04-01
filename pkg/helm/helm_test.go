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

package helm

import (
	"context"
	"testing"

	"github.com/fairwindsops/pluto/v3/pkg/api"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
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
				Component:      "k8s",
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
				Component:      "k8s",
			},
		},
	}
)

func newMockHelm(namespace string) *Helm {
	return &Helm{
		Namespace: namespace,
		Kube:      getMockConfigInstance(),
		Instance: &api.Instance{
			TargetVersions: map[string]string{
				"k8s":          "v1.16.0",
				"istio":        "1.6.1",
				"cert-manager": "v0.15.0",
			},
			DeprecatedVersions: []api.Version{
				{
					Name:           "extensions/v1beta1",
					Kind:           "Deployment",
					DeprecatedIn:   "v1.9.0",
					RemovedIn:      "v1.16.0",
					ReplacementAPI: "apps/v1",
					Component:      "k8s",
				},
				{
					Name:           "apps/v1",
					Kind:           "Deployment",
					DeprecatedIn:   "",
					RemovedIn:      "",
					ReplacementAPI: "",
					Component:      "k8s",
				},
			},
			IgnoreDeprecations: false,
			IgnoreRemovals:     false,
			OutputFormat:       "normal",
		},
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
			want:     nil,
			wantErr:  false,
		},
		{
			name:     "got version",
			manifest: []byte("apiVersion: extensions/v1beta1\nkind: Deployment"),
			want:     []*api.Output{{APIVersion: &api.Version{Name: "extensions/v1beta1", Kind: "Deployment", DeprecatedIn: "v1.9.0", RemovedIn: "v1.16.0", ReplacementAPI: "apps/v1", Component: "k8s"}}},
			wantErr:  false,
		},
		{
			name:     "nil return",
			manifest: []byte("apiVersion: v1beta1\nkind: SomeOtherThing"),
			want:     nil,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newMockHelm("")
			got, err := h.checkForAPIVersion(tt.manifest)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
func TestHelm_getManifestsVersionThree(t *testing.T) {
	tests := []struct {
		name       string
		wantErr    bool
		errMessage string
		secret     *v1.Secret
		want       []*api.Output
	}{
		{
			name:   "helm 3 valid",
			secret: &helmSecret,
			want:   wantOutput,
		},
	}

	for _, tt := range tests {
		h := newMockHelm("")
		if tt.secret != nil {
			ns := v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
			}
			_, err := h.Kube.Client.CoreV1().Namespaces().Create(context.TODO(), &ns, metav1.CreateOptions{})
			if err != nil {
				t.Errorf("failed creating default namespace. test: %s", tt.name)
			}
			_, err = h.Kube.Client.CoreV1().Secrets("default").Create(context.TODO(), tt.secret, metav1.CreateOptions{})
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
			assert.Equal(t, tt.want, h.Instance.Outputs)
		})
	}
}

func TestHelm_getManifest_badClient(t *testing.T) {
	tests := []struct {
		name       string
		wantErr    bool
		errMessage string
		secret     *v1.Secret
		want       []*api.Output
	}{
		{
			name:       "three - bad client",
			wantErr:    true,
			errMessage: "helm 3 function called without helm 3 version set",
		},
	}

	for _, tt := range tests {
		h := &Helm{
			Kube: newBadKubeClient(),
		}
		t.Run(tt.name, func(t *testing.T) {
			err := h.getReleasesVersionThree()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "connect: connection refused")
				return
			}
		})
	}
}

func Test_helmToRelease(t *testing.T) {
	tests := []struct {
		name        string
		helmRelease interface{}
		want        *Release
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "test err in json.Marshal",
			helmRelease: map[string]interface{}{"foo": make(chan int)},
			want:        nil,
			wantErr:     true,
			errMsg:      "error marshaling release: json: unsupported type: chan int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := helmToRelease(tt.helmRelease)
			if tt.wantErr {
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

		})
	}
}
