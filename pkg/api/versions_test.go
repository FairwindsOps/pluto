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

package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/mod/semver"
)

func Test_jsonToStub(t *testing.T) {

	tests := []struct {
		name    string
		data    []byte
		want    []*Stub
		wantErr bool
	}{
		{
			name:    "json not stub",
			data:    []byte("{}"),
			want:    []*Stub{{}},
			wantErr: false,
		},
		{
			name:    "no data",
			data:    []byte(""),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "json is stub",
			data:    []byte(`{"kind": "foo", "apiVersion": "bar"}`),
			want:    []*Stub{{Kind: "foo", APIVersion: "bar"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonToStub(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_yamlToStub(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    []*Stub
		wantErr bool
	}{
		{
			name:    "yaml not stub",
			data:    []byte("foo: bar"),
			want:    []*Stub{{}},
			wantErr: false,
		},
		{
			name:    "not yaml",
			data:    []byte("some text\nthat is not yaml"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "yaml is stub",
			data:    []byte("kind: foo\napiVersion: bar"),
			want:    []*Stub{{Kind: "foo", APIVersion: "bar"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := yamlToStub(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_containsStub(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    []*Stub
		wantErr bool
	}{
		{
			name:    "yaml not stub",
			data:    []byte("foo: bar"),
			want:    []*Stub{{}},
			wantErr: false,
		},
		{
			name:    "not yaml",
			data:    []byte("some text\nthat is not yaml"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "yaml is stub",
			data:    []byte("kind: foo\napiVersion: bar"),
			want:    []*Stub{{Kind: "foo", APIVersion: "bar"}},
			wantErr: false,
		},
		{
			name:    "json not stub",
			data:    []byte("{}"),
			want:    []*Stub{{}},
			wantErr: false,
		},
		{
			name:    "empty string",
			data:    []byte(""),
			want:    nil,
			wantErr: false,
		},
		{
			name:    "no data",
			data:    []byte{},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "json is stub",
			data:    []byte(`{"kind": "foo", "apiVersion": "bar"}`),
			want:    []*Stub{{Kind: "foo", APIVersion: "bar"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := containsStub(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_IsVersioned(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    []*Output
		wantErr bool
	}{
		{
			name:    "yaml no version",
			data:    []byte("foo: bar"),
			want:    nil,
			wantErr: false,
		},
		{
			name:    "not json or yaml",
			data:    []byte("some text\nthat is not yaml"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "yaml has version",
			data:    []byte("kind: Deployment\napiVersion: apps/v1"),
			want:    []*Output{{APIVersion: &Version{Name: "apps/v1", Kind: "Deployment", DeprecatedIn: "", RemovedIn: "", ReplacementAPI: ""}}},
			wantErr: false,
		},
		{
			name:    "json no version",
			data:    []byte("{}"),
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty string",
			data:    []byte(""),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no data",
			data:    []byte{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "json has version",
			data:    []byte(`{"kind": "Deployment", "apiVersion": "extensions/v1beta1"}`),
			want:    []*Output{{APIVersion: &Version{Kind: "Deployment", Name: "extensions/v1beta1", RemovedIn: "v1.16.0", DeprecatedIn: "v1.9.0", ReplacementAPI: "apps/v1"}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsVersioned(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestVersion_IsDeprecatedIn(t *testing.T) {

	tests := []struct {
		name          string
		targetVersion string
		want          bool
		deprecatedIn  string
	}{
		{
			name:          "not deprecated yet 1.15.0",
			targetVersion: "v1.15.0",
			deprecatedIn:  "v1.16.0",
			want:          false,
		},
		{
			name:          "equal values",
			targetVersion: "v1.16.0",
			deprecatedIn:  "v1.16.0",
			want:          true,
		},
		{
			name:          "greater than",
			targetVersion: "v1.17.0",
			deprecatedIn:  "v1.16.0",
			want:          true,
		},
		{
			name:          "Bad semVer",
			targetVersion: "foo",
			deprecatedIn:  "v1.16.0",
			want:          false,
		},
		{
			name:          "blank deprecatedIn - not deprecated",
			targetVersion: "v1.16.0",
			deprecatedIn:  "",
			want:          false,
		},
	}
	for _, tt := range tests {
		deprecatedVersion := &Version{DeprecatedIn: tt.deprecatedIn}
		got := deprecatedVersion.IsDeprecatedIn(tt.targetVersion)
		assert.Equal(t, tt.want, got, "test failed: "+tt.name)
	}
}

func Test_VersionListIsValid(t *testing.T) {
	// This test validates that all of the versions in VersionList are valid semVer
	// it should prevent us from putting bad values in that list in future development
	for _, version := range VersionList {
		if version.DeprecatedIn != "" {
			assert.True(t, semver.IsValid(version.DeprecatedIn), fmt.Sprintf("version %s is not valid semver", version.DeprecatedIn))
		}
		if version.RemovedIn != "" {
			assert.True(t, semver.IsValid(version.RemovedIn), fmt.Sprintf("version %s is not valid semver", version.RemovedIn))
		}
	}
}
