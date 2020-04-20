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
