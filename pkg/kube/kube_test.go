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

package kube

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func Test_getKubeClient(t *testing.T) {
	tests := []struct {
		name           string
		kubeContext    string
		kubeConfig     string
		kubeConfigPath string
		wantErr        bool
	}{
		{
			name:           "context does not exist",
			kubeContext:    "farglebargle",
			kubeConfig:     "testdata/kubeconfig",
			kubeConfigPath: "",
			wantErr:        true,
		},
		{
			name:           "context exists",
			kubeContext:    "kind-kind",
			kubeConfig:     "testdata/kubeconfig",
			kubeConfigPath: "",
			wantErr:        false,
		},
		{
			name:           "invalid kubeconfig",
			kubeContext:    "kind-kind",
			kubeConfig:     "testdata/kubeconfig_invalid",
			kubeConfigPath: "",
			wantErr:        true,
		},
		{
			name:           "invalid kubeconfig",
			kubeContext:    "kind-kind",
			kubeConfig:     "testdata/kubeconfig_invalid",
			kubeConfigPath: "",
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("KUBECONFIG", tt.kubeConfig)
			_, err := GetConfig(tt.kubeContext, tt.kubeConfigPath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
