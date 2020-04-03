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

var testOutput1 = &Output{
	Name: "some name one",
	APIVersion: &Version{
		Name:       "apps/v1",
		Kind:       "Deployment",
		Deprecated: false,
	},
}
var testOutput2 = &Output{
	Name: "some name two",
	APIVersion: &Version{
		Name:       "extensions/v1beta1",
		Kind:       "Deployment",
		Deprecated: true,
	},
}

func ExampleDisplayOutput_showNonDeprecated() {
	_ = DisplayOutput([]*Output{testOutput1}, "tabular", true)

	// Output:
	// KIND         VERSION   DEPRECATED   RESOURCE NAME
	// Deployment   apps/v1   false        some name one
}

func ExampleDisplayOutput() {
	_ = DisplayOutput([]*Output{testOutput1, testOutput2}, "tabular", false)

	// Output:
	// KIND         VERSION              DEPRECATED   RESOURCE NAME
	// Deployment   extensions/v1beta1   true         some name two
}

func ExampleDisplayOutput_json() {
	_ = DisplayOutput([]*Output{testOutput1}, "json", true)

	// Output:
	// [{"file":"some name one","api":{"version":"apps/v1","kind":"Deployment"}}]
}

func ExampleDisplayOutput_yaml() {
	_ = DisplayOutput([]*Output{testOutput1}, "yaml", true)

	// Output:
	// - file: some name one
	//   api:
	//     version: apps/v1
	//     kind: Deployment
}

func ExampleDisplayOutput_noOutput() {
	_ = DisplayOutput([]*Output{testOutput1}, "tabular", false)

	// Output: APIVersions were found, but none were deprecated. Try --show-all.
}

func ExampleDisplayOutput_badFormat() {
	_ = DisplayOutput([]*Output{testOutput1}, "foo", true)

	// Output: output format should be one of (json,yaml,tabular)
}

func ExampleDisplayOutput_zeroLength() {
	_ = DisplayOutput([]*Output{}, "tabular", false)

	// Output: There were no apiVersions found that match our records.
}
