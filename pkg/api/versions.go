package api

import "strings"

// Version is a file that has an apiVersion in it
type Version struct {
	// Name is the name of the api version
	Name string
	// Deprecated is a boolean that indicates if the apiVersion is deprecated
	Deprecated bool
}

// VersionList is a map of apiVersions and if they are deprecated or not.
// TODO: Find a better way to generate this.
var VersionList = []Version{
	{"extensions/v1beta1", true},
	{"apps/v1", false},
}

// StringContainsVersion returns an api version if it is in the string
func StringContainsVersion(line string) *Version {
	for _, version := range VersionList {
		if strings.Contains(line, version.Name) {
			return &version
		}
	}
	return nil
}
