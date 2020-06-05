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

	"github.com/fairwindsops/pluto/pkg/api"
	"github.com/stretchr/testify/assert"
)

var testPath = "testdata"

var deploymentAppsV1Yaml = "testdata/deployment-apps-v1.yaml"
var deploymentAppsV1YamlFile = []*api.Output{{
	Name:      "utilities",
	Namespace: "test-namespaces",
	APIVersion: &api.Version{
		Name:         "apps/v1",
		Kind:         "Deployment",
		DeprecatedIn: "",
		Component:    "k8s",
	}},
}

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
	deploymentAppsV1Yaml,
	deploymentExtensionsV1JSON,
	deploymentExtensionsV1Yaml,
	"testdata/other.txt",
}

var testOutput = []*api.Output{
	deploymentAppsV1YamlFile[0],
	deploymentExtensionsV1JSONFile[0],
	deploymentExtensionsV1YamlFile[0],
}

func TestNewFinder(t *testing.T) {

	wd, _ := os.Getwd()
	tests := []struct {
		name string
		path string
		want *Dir
	}{
		{
			name: "basic",
			path: testPath,
			want: &Dir{RootPath: testPath},
		},
		{
			// This is a bit silly, but it does test the return of the function.
			name: "empty",
			path: "",
			want: &Dir{RootPath: wd},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFinder(tt.path)
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
			dir := NewFinder(tt.directory)
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
			file:    deploymentAppsV1Yaml,
			wantErr: false,
			want:    deploymentAppsV1YamlFile,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckForAPIVersion(tt.file)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
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
			fileList: []string{deploymentAppsV1Yaml},
			want:     deploymentAppsV1YamlFile,
		},
	}
	dir := &Dir{
		RootPath: testPath,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir.FileList = tt.fileList
			err := dir.scanFiles()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.fileList, dir.FileList)
				assert.EqualValues(t, tt.want, dir.Outputs)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := &Dir{
				RootPath: tt.path,
			}
			err := dir.FindVersions()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, dir.Outputs)
			}
		})
	}
}
