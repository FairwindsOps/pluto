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

	"github.com/fairwindsops/pluto/pkg/api"
)

func newMockHelm(version string) *Helm {
	return &Helm{
		Version: version,
		Kube:    getMockConfigInstance(),
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
			name:        "three - error",
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
	}

	for _, tt := range tests {
		h := newMockHelm(tt.helmVersion)
		t.Run(tt.name, func(t *testing.T) {
			err := h.getManifestsVersionThree()
			if tt.wantErr {
				assert.EqualError(t, err, tt.errMessage)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, h.Outputs)
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
