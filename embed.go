// plutoversionsfile makes the Pluto versions.yaml file available to the pluto
// binary and others that import Pluto packages.
package plutoversionsfile

import (
	_ "embed"
)

var (
	//go:embed versions.yaml
	plutoVersionsFileContent []byte
)

// Content returns the Pluto versions.yaml file as a slice of bytes.
func Content() []byte {
	return plutoVersionsFileContent
}
