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

package finder

import (
	"os"
	"testing"

	"github.com/fairwindsops/pluto/v3/pkg/api"
	"github.com/stretchr/testify/assert"
)

var testPath = "testdata"

var deploymentExtensionsV1Yaml = "testdata/deployment-extensions-v1beta1.yaml"
var deploymentExtensionsV1YamlFile = []*api.Output{{
	Name:      "utilities",
	Namespace: "yaml-namespace",
	APIVersion: &api.Version{
		Name:           "extensions/v1beta1",
		Kind:           "Deployment",
		DeprecatedIn:   "v1.9.0",
		RemovedIn:      "v1.16.0",
		ReplacementAPI: "apps/v1",
		Component:      "k8s",
	}},
}

var deploymentExtensionsV1JSON = "testdata/deployment-extensions-v1beta1.json"
var deploymentExtensionsV1JSONFile = []*api.Output{{
	Name:      "utilities",
	Namespace: "json-namespace",
	APIVersion: &api.Version{
		Name:           "extensions/v1beta1",
		Kind:           "Deployment",
		DeprecatedIn:   "v1.9.0",
		RemovedIn:      "v1.16.0",
		ReplacementAPI: "apps/v1",
		Component:      "k8s",
	}},
}

var testFiles = []string{
	deploymentExtensionsV1JSON,
	deploymentExtensionsV1Yaml,
	"testdata/other.txt",
}

var testOutput = []*api.Output{
	deploymentExtensionsV1JSONFile[0],
	deploymentExtensionsV1YamlFile[0],
}

var testVersionDeployment = api.Version{
	Name:           "extensions/v1beta1",
	Kind:           "Deployment",
	DeprecatedIn:   "v1.9.0",
	RemovedIn:      "v1.16.0",
	ReplacementAPI: "apps/v1",
	Component:      "k8s",
}

func newMockFinder(path string) *Dir {
	dir := &Dir{
		RootPath: path,
		Instance: &api.Instance{
			TargetVersions: map[string]string{
				"k8s":          "v1.16.0",
				"istio":        "1.6.1",
				"cert-manager": "v0.15.0",
			},
			DeprecatedVersions: []api.Version{
				testVersionDeployment,
			},
			IgnoreDeprecations: false,
			IgnoreRemovals:     false,
			OutputFormat:       "normal",
		},
	}
	return dir
}

// patchFilePath exists because the filePath will be different
// on every system. This asserts that the current working directory
// is in the file path and then sets it to an empty string
func patchFilePath(t *testing.T, outputs []*api.Output) {
	cwd, _ := os.Getwd()
	// Account for current working dir
	for _, output := range outputs {
		assert.Contains(t, output.FilePath, cwd)
		output.FilePath = ""
	}
}

func TestNewFinder(t *testing.T) {
	wd, _ := os.Getwd()
	tests := []struct {
		name string
		path string
		want *Dir
	}{
		{
			name: "one",
			path: testPath,
			want: newMockFinder(testPath),
		},
		{
			name: "cwd",
			path: "",
			want: &Dir{
				RootPath: wd,
				Instance: &api.Instance{
					TargetVersions: map[string]string{
						"k8s":          "v1.16.0",
						"istio":        "1.6.1",
						"cert-manager": "v0.15.0",
					},
					DeprecatedVersions: []api.Version{
						testVersionDeployment,
					},
					IgnoreDeprecations: false,
					IgnoreRemovals:     false,
					OutputFormat:       "normal",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFinder(tt.path, &api.Instance{
				TargetVersions: map[string]string{
					"k8s":          "v1.16.0",
					"istio":        "1.6.1",
					"cert-manager": "v0.15.0",
				},
				DeprecatedVersions: []api.Version{
					testVersionDeployment,
				},
				IgnoreDeprecations: false,
				IgnoreRemovals:     false,
				OutputFormat:       "normal",
			},
			)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDir_listFiles(t *testing.T) {

	tests := []struct {
		name      string
		wantErr   bool
		fileList  []string
		directory string
	}{
		{
			name:      "pass",
			wantErr:   false,
			fileList:  testFiles,
			directory: testPath,
		},
		{
			name:      "fail",
			wantErr:   true,
			fileList:  []string{},
			directory: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := newMockFinder(testPath)
			dir.RootPath = tt.directory
			err := dir.listFiles()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.fileList, dir.FileList)
			}
		})
	}
}

func Test_checkForAPIVersion(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    []*api.Output
		wantErr bool
	}{
		{
			name:    "deployments extensions/v1beta1",
			file:    deploymentExtensionsV1Yaml,
			wantErr: false,
			want:    deploymentExtensionsV1YamlFile,
		},
		{
			name:    "file dne",
			file:    "foo",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := newMockFinder(testPath)
			got, err := dir.CheckForAPIVersion(tt.file)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				patchFilePath(t, got)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestDir_scanFiles(t *testing.T) {

	tests := []struct {
		name     string
		wantErr  bool
		fileList []string
		want     []*api.Output
	}{
		{
			name:     "pass",
			wantErr:  false,
			fileList: []string{deploymentExtensionsV1Yaml},
			want:     deploymentExtensionsV1YamlFile,
		},
	}
	dir := newMockFinder(testPath)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir.FileList = tt.fileList
			err := dir.scanFiles()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.fileList, dir.FileList)
				patchFilePath(t, dir.Instance.Outputs)
				assert.EqualValues(t, tt.want, dir.Instance.Outputs)
			}
		})
	}
}

func TestDir_FindVersions(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		path    string
		want    []*api.Output
	}{
		{
			name:    "pass",
			wantErr: false,
			path:    testPath,
			want:    testOutput,
		},
		{
			name:    "fail",
			wantErr: true,
			path:    "foo",
			want:    testOutput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := newMockFinder(tt.path)
			dir.Instance.Outputs = nil
			err := dir.FindVersions()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				patchFilePath(t, dir.Instance.Outputs)
				assert.EqualValues(t, tt.want, dir.Instance.Outputs)
			}
		})
	}
}
