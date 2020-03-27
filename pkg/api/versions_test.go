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
	"testing"

	"github.com/stretchr/testify/assert"
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
			want:    []*Output{{APIVersion: &Version{Name: "apps/v1", Kind: "Deployment", Deprecated: false}}},
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
			want:    []*Output{{APIVersion: &Version{Kind: "Deployment", Name: "extensions/v1beta1", Deprecated: true}}},
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
