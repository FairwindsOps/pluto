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

package api

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test_data/versions.yaml
var testVersionsFile []byte

var mockInstance = Instance{
	TargetVersions: map[string]string{
		"k8s":          "v1.16.0",
		"istio":        "1.6.1",
		"cert-manager": "v0.15.0",
	},
	DeprecatedVersions: []Version{
		testVersionDeployment,
	},
}

var testVersionDeploymentString = `deprecated-versions:
- version: extensions/v1beta1
  kind: Deployment
  deprecated-in: v1.9.0
  removed-in: v1.16.0
  replacement-api: apps/v1
  replacement-available-in: v1.10.0
  component: k8s
target-versions:
  k8s: v1.16.0
  istio: v1.6.0
  cert-manager: v0.15.1`

var testVersionDeployment = Version{
	Name:                   "extensions/v1beta1",
	Kind:                   "Deployment",
	DeprecatedIn:           "v1.9.0",
	RemovedIn:              "v1.16.0",
	ReplacementAPI:         "apps/v1",
	ReplacementAvailableIn: "v1.10.0",
	Component:              "k8s",
}

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
		{
			name:    "json list is multiple stubs",
			data:    []byte(`{"kind": "List", "apiVersion": "v1", "items": [{"kind": "foo", "apiVersion": "bar"},{"kind": "bar", "apiVersion": "foo"}]}`),
			want:    []*Stub{{Kind: "foo", APIVersion: "bar"}, {Kind: "bar", APIVersion: "foo"}},
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
			data:    []byte("*."),
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
			name:    "yaml list is multiple stubs",
			data:    []byte("kind: List\napiVersion: v1\nitems:\n- kind: foo\n  apiVersion: bar\n- kind: bar\n  apiVersion: foo"),
			want:    []*Stub{{Kind: "foo", APIVersion: "bar"}, {Kind: "bar", APIVersion: "foo"}},
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
			data:    []byte("*."),
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
			data:    []byte("kind: Deployment\napiVersion: extensions/v1beta1"),
			want:    []*Output{{APIVersion: &testVersionDeployment}},
			wantErr: false,
		},
		{
			name:    "yaml list has version",
			data:    []byte("kind: List\napiVersion: v1\nitems:\n- kind: Deployment\n  apiVersion: extensions/v1beta1"),
			want:    []*Output{{APIVersion: &testVersionDeployment}},
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
			wantErr: false,
		},
		{
			name:    "no data",
			data:    []byte{},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "json has version",
			data:    []byte(`{"kind": "Deployment", "apiVersion": "extensions/v1beta1"}`),
			want:    []*Output{{APIVersion: &testVersionDeployment}},
			wantErr: false,
		},
		{
			name:    "json list has version",
			data:    []byte(`{"kind": "List", "apiVersion": "v1", "items": [{"kind": "Deployment", "apiVersion": "extensions/v1beta1"}]}`),
			want:    []*Output{{APIVersion: &testVersionDeployment}},
			wantErr: false,
		},
		{
			name:    "not yaml",
			data:    []byte("*."),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockInstance.IsVersioned(tt.data)
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
		name           string
		targetVersions map[string]string
		component      string
		want           bool
		deprecatedIn   string
	}{
		{
			name:           "not deprecated yet 1.15.0",
			targetVersions: map[string]string{"foo": "v1.15.0"},
			component:      "foo",
			deprecatedIn:   "v1.16.0",
			want:           false,
		},
		{
			name:           "equal values",
			targetVersions: map[string]string{"foo": "v1.16.0"},
			component:      "foo",
			deprecatedIn:   "v1.16.0",
			want:           true,
		},
		{
			name:           "greater than",
			targetVersions: map[string]string{"foo": "v1.17.0"},
			component:      "foo",
			deprecatedIn:   "v1.16.0",
			want:           true,
		},
		{
			name:           "Bad semVer",
			targetVersions: map[string]string{"foo": "foo"},
			component:      "foo",
			deprecatedIn:   "v1.16.0",
			want:           false,
		},
		{
			name:           "blank deprecatedIn - not deprecated",
			targetVersions: map[string]string{"foo": "v1.16.0"},
			component:      "foo",
			deprecatedIn:   "",
			want:           false,
		},
		{
			name:           "targetversion not included",
			targetVersions: map[string]string{"one": "v1.16.0"},
			component:      "two",
			deprecatedIn:   "v1.16.0",
			want:           false,
		},
	}
	for _, tt := range tests {
		deprecatedVersion := &Version{DeprecatedIn: tt.deprecatedIn, Component: tt.component}
		got := deprecatedVersion.isDeprecatedIn(tt.targetVersions)
		assert.Equal(t, tt.want, got, "test failed: "+tt.name)
	}
}

func TestVersion_IsRemovedIn(t *testing.T) {

	tests := []struct {
		name           string
		targetVersions map[string]string
		component      string
		want           bool
		removedIn      string
	}{
		{
			name:           "not removed yet 1.15.0",
			targetVersions: map[string]string{"foo": "v1.15.0"},
			component:      "foo",
			removedIn:      "v1.16.0",
			want:           false,
		},
		{
			name:           "equal values",
			targetVersions: map[string]string{"foo": "v1.16.0"},
			component:      "foo",
			removedIn:      "v1.16.0",
			want:           true,
		},
		{
			name:           "greater than",
			targetVersions: map[string]string{"foo": "v1.17.0"},
			component:      "foo",
			removedIn:      "v1.16.0",
			want:           true,
		},
		{
			name:           "bad semVer",
			targetVersions: map[string]string{"foo": "foo"},
			removedIn:      "v1.16.0",
			want:           false,
		},
		{
			name:           "blank removedIn - not removed",
			targetVersions: map[string]string{"foo": "v1.16.0"},
			component:      "foo",
			removedIn:      "",
			want:           false,
		},
		{
			name:           "targetVersions not included for component",
			targetVersions: map[string]string{"one": "v1.16.0"},
			component:      "two",
			removedIn:      "v1.16.0",
			want:           false,
		},
	}
	for _, tt := range tests {
		removedVersion := &Version{RemovedIn: tt.removedIn, Component: tt.component}
		got := removedVersion.isRemovedIn(tt.targetVersions)
		assert.Equal(t, tt.want, got, "test failed: "+tt.name)
	}
}

func TestVersion_isReplacementAvailableIn(t *testing.T) {
	tests := []struct {
		name                   string
		targetVersions         map[string]string
		component              string
		want                   bool
		replacementAvailableIn string
	}{
		{
			name:                   "not available yet 1.15.0",
			targetVersions:         map[string]string{"foo": "v1.15.0"},
			component:              "foo",
			replacementAvailableIn: "v1.16.0",
			want:                   false,
		},
		{
			name:                   "equal values",
			targetVersions:         map[string]string{"foo": "v1.16.0"},
			component:              "foo",
			replacementAvailableIn: "v1.16.0",
			want:                   true,
		},
		{
			name:                   "greater than",
			targetVersions:         map[string]string{"foo": "v1.17.0"},
			component:              "foo",
			replacementAvailableIn: "v1.16.0",
			want:                   true,
		},
		{
			name:                   "bad semVer",
			targetVersions:         map[string]string{"foo": "foo"},
			replacementAvailableIn: "v1.16.0",
			want:                   false,
		},
		{
			name:                   "blank replacementAvailableIn - is available",
			targetVersions:         map[string]string{"foo": "v1.16.0"},
			component:              "foo",
			replacementAvailableIn: "",
			want:                   false,
		},
		{
			name:                   "targetVersions not included for component",
			targetVersions:         map[string]string{"one": "v1.16.0"},
			component:              "two",
			replacementAvailableIn: "v1.16.0",
			want:                   false,
		},
	}
	for _, tt := range tests {
		removedVersion := &Version{ReplacementAvailableIn: tt.replacementAvailableIn, Component: tt.component}
		got := removedVersion.isReplacementAvailableIn(tt.targetVersions)
		assert.Equal(t, tt.want, got, "test failed: "+tt.name)
	}
}

func ExampleInstance_printVersionsTabular() {
	instance := Instance{
		DeprecatedVersions: []Version{
			testVersionDeployment,
			{Kind: "testkind", Name: "testname", DeprecatedIn: "", RemovedIn: "", ReplacementAvailableIn: "", Component: "custom"},
		},
	}
	_ = instance.printVersionsTabular()

	// Output:
	// KIND-------- NAME---------------- DEPRECATED IN-- REMOVED IN-- REPLACEMENT-- REPL AVAIL IN-- COMPONENT--
	// Deployment-- extensions/v1beta1-- v1.9.0--------- v1.16.0----- apps/v1------ v1.10.0-------- k8s--------
	// testkind---- testname------------ n/a------------ n/a--------- n/a---------- n/a------------ custom-----
}

func ExampleInstance_printVersionsTabular_noHeaders() {
	instance := Instance{
		DeprecatedVersions: []Version{
			testVersionDeployment,
			{Kind: "testkind", Name: "testname", DeprecatedIn: "", RemovedIn: "", ReplacementAvailableIn: "", Component: "custom"},
		},
		NoHeaders: true,
	}
	_ = instance.printVersionsTabular()

	// Output:
	// Deployment-- extensions/v1beta1-- v1.9.0-- v1.16.0-- apps/v1-- v1.10.0-- k8s-----
	// testkind---- testname------------ n/a----- n/a------ n/a------ n/a------ custom--
}

func ExampleInstance_PrintVersionList_json() {
	instance := Instance{
		DeprecatedVersions: []Version{testVersionDeployment},
	}
	_ = instance.PrintVersionList("json")

	// Output:
	// {"deprecated-versions":[{"version":"extensions/v1beta1","kind":"Deployment","deprecated-in":"v1.9.0","removed-in":"v1.16.0","replacement-api":"apps/v1","replacement-available-in":"v1.10.0","component":"k8s"}]}
}

func ExampleInstance_PrintVersionList_yaml() {
	instance := Instance{
		DeprecatedVersions: []Version{testVersionDeployment},
	}
	_ = instance.PrintVersionList("yaml")

	// Output:
	// deprecated-versions:
	//     - version: extensions/v1beta1
	//       kind: Deployment
	//       deprecated-in: v1.9.0
	//       removed-in: v1.16.0
	//       replacement-api: apps/v1
	//       replacement-available-in: v1.10.0
	//       component: k8s
}

func ExampleInstance_PrintVersionList_normal() {
	instance := Instance{
		DeprecatedVersions: []Version{testVersionDeployment},
	}
	_ = instance.PrintVersionList("normal")

	// Output:
	// KIND-------- NAME---------------- DEPRECATED IN-- REMOVED IN-- REPLACEMENT-- REPL AVAIL IN-- COMPONENT--
	// Deployment-- extensions/v1beta1-- v1.9.0--------- v1.16.0----- apps/v1------ v1.10.0-------- k8s--------
}

func ExampleInstance_PrintVersionList_wide() {
	instance := Instance{
		DeprecatedVersions: []Version{testVersionDeployment},
	}
	_ = instance.PrintVersionList("wide")

	// Output:
	// KIND-------- NAME---------------- DEPRECATED IN-- REMOVED IN-- REPLACEMENT-- REPL AVAIL IN-- COMPONENT--
	// Deployment-- extensions/v1beta1-- v1.9.0--------- v1.16.0----- apps/v1------ v1.10.0-------- k8s--------
}

func ExampleInstance_PrintVersionList_badformat() {
	instance := Instance{
		DeprecatedVersions: []Version{testVersionDeployment},
	}
	_ = instance.PrintVersionList("foo")

	// Output:
	// The output format must be one of (normal|wide|json|yaml)
}

func Test_isDuplicate(t *testing.T) {

	tests := []struct {
		name string
		a    Version
		b    Version
		want bool
	}{
		{
			name: "is duplicate",
			a:    Version{Kind: "Deployment", Name: "apps/v1"},
			b:    Version{Kind: "Deployment", Name: "apps/v1"},
			want: true,
		},
		{
			name: "is not duplicate",
			a:    Version{Kind: "Deployment", Name: "extensions/v1beta1"},
			b:    Version{Kind: "Deployment", Name: "apps/v1"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isDuplicate(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVersion_isContainedIn(t *testing.T) {
	tests := []struct {
		name        string
		version     Version
		versionList []Version
		want        bool
	}{
		{
			name:    "true",
			version: Version{Kind: "Deployment", Name: "extensions/v1beta1"},
			versionList: []Version{
				{Kind: "Deployment", Name: "extensions/v1beta1"},
				{Kind: "Deployment", Name: "apps/v1"},
			},
			want: true,
		},
		{
			name:    "false",
			version: Version{Kind: "Deployment", Name: "extensions/v1beta1"},
			versionList: []Version{
				{Kind: "Deployment", Name: "apps/v1"},
				{Kind: "Deployment", Name: "extensions/v1beta2"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.version.isContainedIn(tt.versionList)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCombineAdditionalVersions(t *testing.T) {
	type args struct {
		additional []Version
		defaults   []Version
	}
	tests := []struct {
		name     string
		args     args
		want     []Version
		wantErr  bool
		errorMsg string
	}{
		{
			name: "error combining due to duplicate",
			args: args{
				additional: []Version{
					{Kind: "Deployment", Name: "apps/v1"},
				},
				defaults: []Version{
					{Kind: "Deployment", Name: "apps/v1"},
				},
			},
			wantErr:  true,
			errorMsg: "duplicate cannot be added to defaults: Deployment apps/v1",
		},
		{
			name: "error combining due to duplicate",
			args: args{
				additional: []Version{
					{Kind: "Deployment", Name: "extensions/v1beta1"},
				},
				defaults: []Version{
					{Kind: "Deployment", Name: "apps/v1"},
				},
			},
			wantErr: false,
			want: []Version{
				{Kind: "Deployment", Name: "apps/v1"},
				{Kind: "Deployment", Name: "extensions/v1beta1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CombineAdditionalVersions(tt.args.additional, tt.args.defaults)
			if tt.wantErr {
				assert.EqualError(t, err, tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestMarshalVersions(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    []Version
		wantErr bool
	}{
		{
			name:    "no error",
			data:    []byte(testVersionDeploymentString),
			want:    []Version{testVersionDeployment},
			wantErr: false,
		},
		{
			name:    "unmarshal error",
			data:    []byte(`foo`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := UnMarshalVersions(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestGetDefaultVersionList(t *testing.T) {

	// This test will ensure that the versions.yaml file is well-formed and doesn't break anything.
	defaultVersions, defaultTargetVersions, err := GetDefaultVersionList(testVersionsFile)
	assert.NoError(t, err)
	assert.NotNil(t, defaultVersions)
	assert.NotNil(t, defaultTargetVersions)
}

func TestInstance_checkVersion(t *testing.T) {
	tests := []struct {
		name     string
		instance *Instance
		stub     *Stub
		want     *Version
	}{
		{
			name: "empty kind",
			instance: &Instance{
				DeprecatedVersions: []Version{
					{Kind: "", Name: "cert-manager.k8s.io", Component: "cert-manager"},
				},
			},
			stub: &Stub{
				Kind:       "any",
				APIVersion: "cert-manager.k8s.io",
				Metadata: StubMeta{
					Name:      "foo",
					Namespace: "foobar",
				},
			},
			want: &Version{
				Name:      "cert-manager.k8s.io",
				Kind:      "any",
				Component: "cert-manager",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.instance.checkVersion(tt.stub)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
