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

var testOutput1 = &Output{
	Name: "some name one",
	APIVersion: &Version{
		Name:           "apps/v1",
		Kind:           "Deployment",
		DeprecatedIn:   "",
		RemovedIn:      "",
		ReplacementAPI: "",
	},
}
var testOutput2 = &Output{
	Name: "some name two",
	APIVersion: &Version{
		Name:           "extensions/v1beta1",
		Kind:           "Deployment",
		DeprecatedIn:   "v1.9.0",
		RemovedIn:      "v1.16.0",
		ReplacementAPI: "apps/v1",
	},
}

func init() {
	padChar = byte('-')
}

func ExampleInstance_DisplayOutput_showAll_normal() {
	instance := &Instance{
		TargetVersion: "v1.15.0",
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		OutputFormat: "normal",
		ShowAll:      true,
	}
	_ = instance.DisplayOutput()

	// Output:
	// NAME----------- KIND-------- VERSION------------- REPLACEMENT-- REMOVED-- DEPRECATED--
	// some name one-- Deployment-- apps/v1------------- ------------- false---- false-------
	// some name two-- Deployment-- extensions/v1beta1-- apps/v1------ false---- true--------
}

func ExampleInstance_DisplayOutput_showAll_wide() {
	instance := &Instance{
		TargetVersion: "v1.16.0",
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		OutputFormat: "wide",
		ShowAll:      true,
	}
	_ = instance.DisplayOutput()

	// Output:
	// NAME----------- KIND-------- VERSION------------- REPLACEMENT-- DEPRECATED-- DEPRECATED IN-- REMOVED-- REMOVED IN--
	// some name one-- Deployment-- apps/v1------------- ------------- false------- --------------- false---- ------------
	// some name two-- Deployment-- extensions/v1beta1-- apps/v1------ true-------- v1.9.0--------- true----- v1.16.0-----
}

func ExampleInstance_DisplayOutput_normal() {
	instance := &Instance{
		TargetVersion: "v1.16.0",
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		OutputFormat: "normal",
		ShowAll:      false,
	}
	_ = instance.DisplayOutput()

	// Output:
	// NAME----------- KIND-------- VERSION------------- REPLACEMENT-- REMOVED-- DEPRECATED--
	// some name two-- Deployment-- extensions/v1beta1-- apps/v1------ true----- true--------
}

func ExampleInstance_DisplayOutput_wide() {
	instance := &Instance{
		TargetVersion: "v1.16.0",
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		OutputFormat: "wide",
		ShowAll:      false,
	}
	_ = instance.DisplayOutput()

	// Output:
	// NAME----------- KIND-------- VERSION------------- REPLACEMENT-- DEPRECATED-- DEPRECATED IN-- REMOVED-- REMOVED IN--
	// some name two-- Deployment-- extensions/v1beta1-- apps/v1------ true-------- v1.9.0--------- true----- v1.16.0-----
}

func ExampleInstance_DisplayOutput_showAll_json() {
	instance := &Instance{
		TargetVersion: "v1.16.0",
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		OutputFormat: "json",
		ShowAll:      true,
	}
	_ = instance.DisplayOutput()

	// Output:
	// {"items":[{"name":"some name one","api":{"version":"apps/v1","kind":"Deployment","deprecated-in":"","removed-in":"","replacement-api":""},"deprecated":false,"removed":false},{"name":"some name two","api":{"version":"extensions/v1beta1","kind":"Deployment","deprecated-in":"v1.9.0","removed-in":"v1.16.0","replacement-api":"apps/v1"},"deprecated":true,"removed":true}],"show-all":true,"target-version":"v1.16.0"}
}

func ExampleInstance_DisplayOutput_showAll_yaml() {
	instance := &Instance{
		TargetVersion: "v1.16.0",
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		OutputFormat: "yaml",
		ShowAll:      true,
	}
	_ = instance.DisplayOutput()

	// Output:
	// items:
	// - name: some name one
	//   api:
	//     version: apps/v1
	//     kind: Deployment
	//     deprecated-in: ""
	//     removed-in: ""
	//     replacement-api: ""
	//   deprecated: false
	//   removed: false
	// - name: some name two
	//   api:
	//     version: extensions/v1beta1
	//     kind: Deployment
	//     deprecated-in: v1.9.0
	//     removed-in: v1.16.0
	//     replacement-api: apps/v1
	//   deprecated: true
	//   removed: true
	// show-all: true
	// target-version: v1.16.0
}

func ExampleInstance_DisplayOutput_noOutput() {
	instance := &Instance{
		TargetVersion: "v1.16.0",
		Outputs: []*Output{
			testOutput1,
		},
		OutputFormat: "normal",
		ShowAll:      false,
	}
	_ = instance.DisplayOutput()

	// Output: APIVersions were found, but none were deprecated. Try --show-all.
}

func ExampleInstance_DisplayOutput_badFormat() {
	instance := &Instance{
		TargetVersion: "v1.16.0",
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		OutputFormat: "foo",
		ShowAll:      false,
	}
	_ = instance.DisplayOutput()

	// Output: output format should be one of (json,yaml,normal,wide)
}

func ExampleInstance_DisplayOutput_zeroLength() {
	instance := &Instance{
		TargetVersion: "v1.16.0",
		Outputs:       []*Output{},
		OutputFormat:  "normal",
		ShowAll:       false,
	}
	_ = instance.DisplayOutput()

	// Output: There were no apiVersions found that match our records.
}

func TestGetReturnCode(t *testing.T) {

	type args struct {
		outputs            []*Output
		ignoreDeprecations bool
		ignoreRemovals     bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "empty return zero",
			args: args{
				outputs:            []*Output{},
				ignoreDeprecations: false,
			},
			want: 0,
		},
		{
			name: "version is deprecated return one",
			args: args{
				outputs: []*Output{
					{
						APIVersion: &Version{
							DeprecatedIn: "v1.16.0",
							RemovedIn:    "v1.20.0",
						},
					},
				},
				ignoreDeprecations: false,
				ignoreRemovals:     false,
			},
			want: 2,
		},
		{
			name: "version is deprecated ignore deprecations",
			args: args{
				outputs: []*Output{
					{
						APIVersion: &Version{
							DeprecatedIn: "v1.16.0",
							RemovedIn:    "v1.20.0",
						},
					},
				},
				ignoreDeprecations: true,
				ignoreRemovals:     false,
			},
			want: 0,
		},
		{
			name: "version is removed",
			args: args{
				outputs: []*Output{
					{
						APIVersion: &Version{
							RemovedIn:    "v1.16.0",
							DeprecatedIn: "v1.12.0",
						},
					},
				},
				ignoreDeprecations: false,
				ignoreRemovals:     false,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance := &Instance{
				TargetVersion:      "v1.16.0",
				IgnoreDeprecations: tt.args.ignoreDeprecations,
				IgnoreRemovals:     tt.args.ignoreRemovals,
				Outputs:            tt.args.outputs,
			}
			got := instance.GetReturnCode()
			assert.Equal(t, tt.want, got)
		})
	}
}
