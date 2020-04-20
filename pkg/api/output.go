package api

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"gopkg.in/yaml.v2"
)

var padChar = byte(' ')

// Output is a thing that has an apiVersion in it
type Output struct {
	Name       string   `json:"file,omitempty" yaml:"file,omitempty"`
	APIVersion *Version `json:"api,omitempty" yaml:"api,omitempty"`
}

// Instance is an instance of the API. This holds configuration for a "run" of Pluto
type Instance struct {
	Outputs            []*Output
	IgnoreDeprecations bool
	IgnoreRemovals     bool
	OutputFormat       string
	ShowAll            bool
	TargetVersion      string
}

// DisplayOutput prints the output based on desired variables
func (instance *Instance) DisplayOutput() error {
	if len(instance.Outputs) == 0 {
		fmt.Println("There were no apiVersions found that match our records.")
		return nil
	}
	instance.filterOutput()
	var err error
	var outData []byte
	switch instance.OutputFormat {
	case "normal":
		t, err := instance.tabOut()
		if err != nil {
			return err
		}
		err = t.Flush()
		if err != nil {
			return err
		}
		return nil
	case "wide":
		t, err := instance.tabOut()
		if err != nil {
			return err
		}
		err = t.Flush()
		if err != nil {
			return err
		}
		return nil
	case "json":
		outData, err = json.Marshal(instance.Outputs)
		if err != nil {
			return err
		}
		fmt.Println(string(outData))
	case "yaml":
		outData, err = yaml.Marshal(instance.Outputs)
		if err != nil {
			return err
		}
		fmt.Println(string(outData))
	default:
		fmt.Println("output format should be one of (json,yaml,normal,wide)")
	}
	return nil
}

func (instance *Instance) filterOutput() {
	if instance.ShowAll {
		return
	}
	var usableOutputs []*Output
	for _, output := range instance.Outputs {
		if output.APIVersion.IsDeprecatedIn(instance.TargetVersion) || output.APIVersion.IsRemovedIn(instance.TargetVersion) {
			usableOutputs = append(usableOutputs, output)
		}
	}
	instance.Outputs = usableOutputs

}

func (instance *Instance) tabOut() (*tabwriter.Writer, error) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 15, 2, padChar, 0)

	if len(instance.Outputs) == 0 {
		_, err := fmt.Fprintln(w, "APIVersions were found, but none were deprecated. Try --show-all.")
		if err != nil {
			return nil, err
		}
		return w, nil
	}

	if instance.OutputFormat == "normal" {
		_, err := fmt.Fprintln(w, "NAME\t KIND\t VERSION\t REPLACEMENT\t REMOVED\t DEPRECATED\t")
		if err != nil {
			return nil, err
		}
		for _, output := range instance.Outputs {
			kind := output.APIVersion.Kind
			removed := fmt.Sprintf("%t", output.APIVersion.IsRemovedIn(instance.TargetVersion))
			deprecated := fmt.Sprintf("%t", output.APIVersion.IsDeprecatedIn(instance.TargetVersion))
			version := output.APIVersion.Name
			name := output.Name
			replacement := output.APIVersion.ReplacementAPI

			_, err = fmt.Fprintf(w, "%s\t %s\t %s\t %s\t %s\t %s\t\n", name, kind, version, replacement, removed, deprecated)
			if err != nil {
				return nil, err
			}
		}
	}

	if instance.OutputFormat == "wide" {
		_, err := fmt.Fprintln(w, "NAME\t KIND\t VERSION\t REPLACEMENT\t DEPRECATED\t DEPRECATED IN\t REMOVED\t REMOVED IN\t")
		if err != nil {
			return nil, err
		}
		for _, output := range instance.Outputs {
			kind := output.APIVersion.Kind
			deprecated := fmt.Sprintf("%t", output.APIVersion.IsDeprecatedIn(instance.TargetVersion))
			removed := fmt.Sprintf("%t", output.APIVersion.IsRemovedIn(instance.TargetVersion))
			version := output.APIVersion.Name
			name := output.Name
			replacement := output.APIVersion.ReplacementAPI
			deprecatedIn := output.APIVersion.DeprecatedIn
			removedIn := output.APIVersion.RemovedIn

			_, err = fmt.Fprintf(w, "%s\t %s\t %s\t %s\t %s\t %s\t %s\t %s\t\n", name, kind, version, replacement, deprecated, deprecatedIn, removed, removedIn)
			if err != nil {
				return nil, err
			}
		}

	}
	return w, nil
}

// GetReturnCode checks for deprecated versions and returns a code.
// takes a boolean to ignore any errors.
// exit 2 - version deprecated
// exit 3 - version removed
func (instance *Instance) GetReturnCode() int {
	returnCode := 0
	var deprecations int
	var removals int
	for _, output := range instance.Outputs {
		if output.APIVersion.IsRemovedIn(instance.TargetVersion) {

			removals = removals + 1
		}
		if output.APIVersion.IsDeprecatedIn(instance.TargetVersion) {
			deprecations = deprecations + 1
		}
	}
	if deprecations > 0 && !instance.IgnoreDeprecations {
		returnCode = 2
	}
	if removals > 0 && !instance.IgnoreRemovals {
		returnCode = 3
	}
	return returnCode
}
