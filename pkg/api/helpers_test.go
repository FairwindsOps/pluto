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

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StringInSlice(t *testing.T) {
	type args struct {
		a    string
		list []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test true",
			args: args{
				a:    "string",
				list: []string{"test", "string"},
			},
			want: true,
		},
		{
			name: "test false",
			args: args{
				a:    "string",
				list: []string{"test", "nothere"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringInSlice(tt.args.a, tt.args.list); got != tt.want {
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestIsFileOrStdin(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"stdin", "-", true},
		{"no file", "notafile.foo", false},
		{"file", "helpers.go", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsFileOrStdin(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
