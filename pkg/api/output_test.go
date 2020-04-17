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
		Name:         "apps/v1",
		Kind:         "Deployment",
		DeprecatedIn: "",
	},
}
var testOutput2 = &Output{
	Name: "some name two",
	APIVersion: &Version{
		Name:         "extensions/v1beta1",
		Kind:         "Deployment",
		DeprecatedIn: "v1.16.0",
	},
}

var targetVersion116 = string("v1.16.0")

func ExampleDisplayOutput_showNonDeprecated() {
	_ = DisplayOutput([]*Output{testOutput1}, "tabular", true, targetVersion116, '-')

	// Output:
	// KIND-------- VERSION-- DEPRECATED-- DEPRECATED IN-- RESOURCE NAME--
	// Deployment-- apps/v1-- false------- n/a------------ some name one--
}

func ExampleDisplayOutput() {
	_ = DisplayOutput([]*Output{testOutput1, testOutput2}, "tabular", false, targetVersion116, '-')

	// Output:
	// KIND-------- VERSION------------- DEPRECATED-- DEPRECATED IN-- RESOURCE NAME--
	// Deployment-- extensions/v1beta1-- true-------- v1.16.0-------- some name two--
}

func ExampleDisplayOutput_json() {
	_ = DisplayOutput([]*Output{testOutput1}, "json", true, targetVersion116, '-')

	// Output:
	// [{"file":"some name one","api":{"version":"apps/v1","kind":"Deployment"}}]
}

func ExampleDisplayOutput_yaml() {
	_ = DisplayOutput([]*Output{testOutput1}, "yaml", true, targetVersion116, '-')

	// Output:
	// - file: some name one
	//   api:
	//     version: apps/v1
	//     kind: Deployment
}

func ExampleDisplayOutput_noOutput() {
	_ = DisplayOutput([]*Output{testOutput1}, "tabular", false, targetVersion116, '-')

	// Output: APIVersions were found, but none were deprecated. Try --show-all.
}

func ExampleDisplayOutput_badFormat() {
	_ = DisplayOutput([]*Output{testOutput1}, "foo", true, targetVersion116, '-')

	// Output: output format should be one of (json,yaml,tabular)
}

func ExampleDisplayOutput_zeroLength() {
	_ = DisplayOutput([]*Output{}, "tabular", false, targetVersion116, '-')

	// Output: There were no apiVersions found that match our records.
}

func TestGetReturnCode(t *testing.T) {
	type args struct {
		outputs      []*Output
		ignoreErrors bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "empty return zero",
			args: args{
				outputs:      []*Output{},
				ignoreErrors: false,
			},
			want: 0,
		},
		{
			name: "version is deprecated return one",
			args: args{
				outputs: []*Output{
					{
						APIVersion: &Version{
							DeprecatedIn: targetVersion116,
						},
					},
				},
				ignoreErrors: false,
			},
			want: 1,
		},
		{
			name: "version is deprecated ignore errors",
			args: args{
				outputs: []*Output{
					{
						APIVersion: &Version{
							DeprecatedIn: targetVersion116,
						},
					},
				},
				ignoreErrors: true,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetReturnCode(tt.args.outputs, tt.args.ignoreErrors, targetVersion116)
			assert.Equal(t, tt.want, got)
		})
	}
}
