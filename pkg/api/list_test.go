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
