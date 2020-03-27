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
	"flag"
	"testing"

	"github.com/fairwindsops/pluto/pkg/api"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	multipleReleaseVersionsList = []*Release{
		{
			Name: "sh.helm.release.v1.example-app.v1",
		},
		{
			Name: "sh.helm.release.v1.example-app.v2",
		},
		{
			Name: "sh.helm.release.v1.example-app.v123",
		},
		{
			Name: "sh.helm.release.v1.example-app.v3",
		},
	}
	badReleaseName = []*Release{
		{
			Name: "sh.helm.release.v1.example-app.vOne",
		},
	}
	gzippedReleaseDataString    = "H4sIACgjfl4AA3WSTW8TMRCG7/kVo0XiREjSVgi5p4qPA0oBIcGF5TDxTpIBf2w9dsIK8d+xt02zm4a5WJ535rHntf9MACqLjtcksVJQYcvfKAh7pwDbVma7Re1+sWsUvKXW+M6Si7WzFLHBiKp2AA4tKUiRDUcmKSmDKzLSq1A4I1la0r0UMpE1ioJF2QoZ0tGHhzaLUW+XQ9BTFEAk2xqMdGgaXqyEGQPOIfLJhxuV0N5FZJddeExNz814H2xxk5W7hN1L9jNJDYVs6M8weyxW5X4Sjz3P4EOSmE9lB89hjxxh7QPtKBxrtLcZk23/DnU1W7GbrVC2dfUib6f6YZ3WFfwYjBY2ct+w37IhiCHRNTQexBC1cDkvG0fX4y4hnQLH7k2em37HgVXlhbD55Ez3xfv4PiOlk2y46snDOjTG7z8H3uWaDb0TjXnm/hOt0cioNCR3Ix+9K8inoF79KhTyn5hfXlwNNY0trg6eDgWAJvj2JFVe7Wa5PCYDiU9Bn/QGukv5cU6Buk0qG2bHWUvWh07Bq6tbHiqGLZ9HLOb/YywuXt9yNfk7+QdmzTSnhAMAAA=="
	notGzippedReleaseDataString = "ewogICJtYW5pZmVzdCI6ICJhcGlWZXJzaW9uOiBhcHBzL3YxXG5raW5kOiBEZXBsb3ltZW50XG5tZXRhZGF0YTpcbiAgbmFtZTogdXRpbGl0aWVzXG4gIGxhYmVsczpcbiAgICBhcHA6IHV0aWxpdGllc1xuc3BlYzpcbiAgcmVwbGljYXM6IDFcbiAgc2VsZWN0b3I6XG4gICAgbWF0Y2hMYWJlbHM6XG4gICAgICBhcHA6IHV0aWxpdGllc1xuICB0ZW1wbGF0ZTpcbiAgICBtZXRhZGF0YTpcbiAgICAgIGxhYmVsczpcbiAgICAgICAgYXBwOiB1dGlsaXRpZXNcbiAgICBzcGVjOlxuICAgICAgY29udGFpbmVyczpcbiAgICAgIC0gbmFtZTogdXRpbGl0aWVzXG4gICAgICAgIGltYWdlOiBxdWF5LmlvL3N1ZGVybWFuanIvdXRpbGl0aWVzOmxhdGVzdFxuICAgICAgICAjIEp1c3Qgc3BpbiAmIHdhaXQgZm9yZXZlclxuICAgICAgICBjb21tYW5kOiBbIFwiL2Jpbi9iYXNoXCIsIFwiLWNcIiwgXCItLVwiIF1cbiAgICAgICAgYXJnczogWyBcIndoaWxlIHRydWU7IGRvIHNsZWVwIDMwOyBkb25lO1wiIF1cbiAgICAgICAgc2VjdXJpdHlDb250ZXh0OlxuICAgICAgICAgIHJlYWRPbmx5Um9vdEZpbGVzeXN0ZW06IHRydWVcbiAgICAgICAgICBhbGxvd1ByaXZpbGVnZUVzY2FsYXRpb246IGZhbHNlXG4gICAgICAgICAgcnVuQXNOb25Sb290OiB0cnVlXG4gICAgICAgICAgcnVuQXNVc2VyOiAxMDMyNFxuICAgICAgICAgIGNhcGFiaWxpdGllczpcbiAgICAgICAgICAgIGRyb3A6XG4gICAgICAgICAgICAgIC0gQUxMXG4gICAgICAgIHJlc291cmNlczpcbiAgICAgICAgICByZXF1ZXN0czpcbiAgICAgICAgICAgIGNwdTogMzBtXG4gICAgICAgICAgICBtZW1vcnk6IDY0TWlcbiAgICAgICAgICBsaW1pdHM6XG4gICAgICAgICAgICBjcHU6IDEwMG1cbiAgICAgICAgICAgIG1lbW9yeTogMTI4TWkiCn0K"
	decodedRelease              = &Release{
		Manifest: "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: utilities\n  labels:\n    app: utilities\nspec:\n  replicas: 1\n  selector:\n    matchLabels:\n      app: utilities\n  template:\n    metadata:\n      labels:\n        app: utilities\n    spec:\n      containers:\n      - name: utilities\n        image: quay.io/sudermanjr/utilities:latest\n        # Just spin & wait forever\n        command: [ \"/bin/bash\", \"-c\", \"--\" ]\n        args: [ \"while true; do sleep 30; done;\" ]\n        securityContext:\n          readOnlyRootFilesystem: true\n          allowPrivilegeEscalation: false\n          runAsNonRoot: true\n          runAsUser: 10324\n          capabilities:\n            drop:\n              - ALL\n        resources:\n          requests:\n            cpu: 30m\n            memory: 64Mi\n          limits:\n            cpu: 100m\n            memory: 128Mi",
	}
	kubeSecret = &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "utilities",
			Namespace: "test-namespace",
		},
		Data: map[string][]byte{
			"release": []byte(gzippedReleaseDataString),
		},
	}
)

func newMockHelm(version string) *Helm {
	return &Helm{
		Version: version,
		Kube:    getMockConfigInstance(),
	}
}

func TestSetCurrentRelease(t *testing.T) {
	tests := []struct {
		name        string
		releaseList []*Release
		want        []*Release
		wantErr     bool
	}{
		{
			name:        "empty",
			releaseList: []*Release{},
			want:        nil,
			wantErr:     false,
		},
		{
			name:        "multiple release versions",
			releaseList: multipleReleaseVersionsList,
			want:        []*Release{multipleReleaseVersionsList[2]},
			wantErr:     false,
		},
		{
			name:        "bad name",
			releaseList: badReleaseName,
			want:        nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Helm{}
			err := h.setCurrentReleases(tt.releaseList)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, h.CurrentReleases)
			}
		})
	}
}

func TestDecodeReleaseSecret(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    *Release
		wantErr bool
	}{
		{
			name:    "empty",
			data:    "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid b64",
			data:    "notb64",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "regular decode",
			data:    gzippedReleaseDataString,
			want:    decodedRelease,
			wantErr: false,
		},
		{
			name:    "regular decode",
			data:    notGzippedReleaseDataString,
			want:    decodedRelease,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeReleaseSecret(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
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
			want:     []*api.Output{{APIVersion: &api.Version{Name: "apps/v1", Kind: "Deployment", Deprecated: false}}},
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
		wantErr     bool
		errMessage  string
	}{
		{
			name:        "basic",
			helmVersion: "2",
			wantErr:     true,
			errMessage:  "helm 2 check not implemented",
		},
		{
			name:        "three",
			helmVersion: "3",
			wantErr:     true,
			errMessage:  "helm 2 function called without helm 2 version set",
		},
	}
	for _, tt := range tests {
		h := newMockHelm(tt.helmVersion)
		t.Run(tt.name, func(t *testing.T) {
			err := h.getManifestsVersionTwo()
			if tt.wantErr {
				assert.EqualError(t, err, tt.errMessage)
				return
			}
			assert.NoError(t, err)

		})
	}
}

func TestHelm_getHelmSecrets(t *testing.T) {
	err := flag.Set("test.v", "false")
	if err != nil {
		t.Fatal(err)
	}
	kubeSecretNonHelm := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "non-helm-secret",
			Namespace: "test-namespace",
		},
		Type: "generic",
	}
	kubeSecretHelm := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "helm-secret",
			Namespace: "test-namespace",
		},
		Type: "helm.sh/release.v1",
	}
	tests := []struct {
		name        string
		helmVersion string
		want        []*v1.Secret
		mockSecrets []*v1.Secret
		wantErr     bool
	}{
		{
			name:        "basic",
			helmVersion: "2",
			want:        []*v1.Secret{},
			wantErr:     true,
		},
		{
			name:        "empty",
			helmVersion: "3",
			want:        []*v1.Secret{&kubeSecretHelm},
			mockSecrets: []*v1.Secret{&kubeSecretNonHelm, &kubeSecretHelm},
			wantErr:     false,
		},
	}
	for _, tt := range tests {

		// Initalize the mock
		h := newMockHelm(tt.helmVersion)
		if len(tt.want) > 0 {
			for _, secret := range tt.mockSecrets {
				_, err := h.Kube.Client.CoreV1().Secrets("test-namespace").Create(secret)
				if err != nil {
					t.Fatal(err)
				}
			}
		}

		// Run the test
		t.Run(tt.name, func(t *testing.T) {
			secrets, err := h.getHelmSecrets()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, tt.want, secrets)
		})
	}
}

func TestHelm_getReleases(t *testing.T) {
	tests := []struct {
		name    string
		secrets []*v1.Secret
		want    []*Release
		wantErr bool
	}{
		{
			name:    "functional",
			secrets: []*v1.Secret{kubeSecret},
			want:    []*Release{&Release{Name: "utilities", Manifest: decodedRelease.Manifest}},
			wantErr: false,
		},
		{
			name:    "empty",
			secrets: []*v1.Secret{},
			want:    nil,
			wantErr: false,
		},
		{
			name: "error",
			secrets: []*v1.Secret{&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "bad",
				},
				Data: map[string][]byte{
					"release": []byte("this is not good data"),
				},
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		h := newMockHelm("3")
		t.Run(tt.name, func(t *testing.T) {
			got, err := h.getReleases(tt.secrets)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, tt.want, got)
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
			h := newMockHelm(tt.helmVersion)
			err := h.FindVersions()
			if tt.wantErr {
				assert.EqualError(t, err, tt.errMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
