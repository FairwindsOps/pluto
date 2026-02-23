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
	"os"
	"slices"
)

// IsFileOrStdin detects if a file exists, or returns true if - is passed
func IsFileOrStdin(name string) bool {
	if name == "-" {
		return true
	}
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// StringInSlice returns true if the string is contained in the slice
func StringInSlice(s string, slice []string) bool {
	return slices.Contains(slice, s)
}
