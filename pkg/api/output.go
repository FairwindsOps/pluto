package api

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"gopkg.in/yaml.v2"
)

var padChar = byte(' ')

// Output is a thing that has an apiVersion in it
type Output struct {
	// Name is the name of the object in question.
	// This might be an object name, or a release
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Namespace is the namespace that the object is in
	// The output may resolve this to UNKNOWN if there is no way of determining it
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	// APIVersion is the version object corresponding to this output
	APIVersion *Version `json:"api,omitempty" yaml:"api,omitempty"`
	// Deprecated is a boolean indicating whether or not the version is deprecated
	Deprecated bool `json:"deprecated" yaml:"deprecated"`
	// Removed is a boolean indicating whether or not the version has been removed
	Removed bool `json:"removed" yaml:"removed"`
	// CustomColumns is a list of column headers you want to show
}

// Instance is an instance of the API. This holds configuration for a "run" of Pluto
type Instance struct {
	Outputs            []*Output         `json:"items,omitempty" yaml:"items,omitempty"`
	IgnoreDeprecations bool              `json:"-" yaml:"-"`
	IgnoreRemovals     bool              `json:"-" yaml:"-"`
	OutputFormat       string            `json:"-" yaml:"-"`
	TargetVersions     map[string]string `json:"target-versions,omitempty" yaml:"target-versions,omitempty"`
	DeprecatedVersions []Version         `json:"-" yaml:"-"`
	CustomColumns      []string          `json:"-" yaml:"-"`
}

// DisplayOutput prints the output based on desired variables
func (instance *Instance) DisplayOutput() error {
	if len(instance.Outputs) == 0 && (instance.OutputFormat == "normal" || instance.OutputFormat == "wide") {
		fmt.Println("There were no resources found with known deprecated apiVersions.")
		return nil
	}
	instance.filterOutput()
	var err error
	var outData []byte
	switch instance.OutputFormat {
	case "normal":
		c := instance.normalColumns()
		t := instance.tabOut(c)
		err = t.Flush()
		if err != nil {
			return err
		}
		return nil
	case "wide":
		c := instance.wideColumns()
		t := instance.tabOut(c)
		err = t.Flush()
		if err != nil {
			return err
		}
		return nil
	case "custom-columns":
		c := instance.customColumns()
		t := instance.tabOut(c)
		err = t.Flush()
		if err != nil {
			return err
		}
	case "json":
		outData, err = json.Marshal(instance)
		if err != nil {
			return err
		}
		fmt.Println(string(outData))
	case "yaml":
		outData, err = yaml.Marshal(instance)
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
	var usableOutputs []*Output
	for _, output := range instance.Outputs {
		output.Deprecated = output.APIVersion.isDeprecatedIn(instance.TargetVersions)
		output.Removed = output.APIVersion.isRemovedIn(instance.TargetVersions)

		if output.Deprecated || output.Removed {
			usableOutputs = append(usableOutputs, output)
		}
	}
	instance.Outputs = usableOutputs

}

func (instance *Instance) tabOut(columns columnList) *tabwriter.Writer {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 15, 2, padChar, 0)

	if len(instance.Outputs) == 0 {
		_, _ = fmt.Fprintln(w, "No output to display")
		return w
	}

	columnIndexes := make([]int, 0, len(columns))
	for k := range columns {
		columnIndexes = append(columnIndexes, k)
	}
	sort.Ints(columnIndexes)

	var headers string
	for _, k := range columnIndexes {
		if k == 0 {
			headers = fmt.Sprintf("%s\t", columns[k].header())
		} else {
			headers = fmt.Sprintf("%s %s\t", headers, columns[k].header())
		}
	}
	_, _ = fmt.Fprintln(w, headers)

	var rows string
	for _, o := range instance.Outputs {
		var row string
		for _, k := range columnIndexes {
			if k == 0 {
				row = fmt.Sprintf("%s\t", columns[k].value(o))
			} else {
				row = fmt.Sprintf("%s %s\t", row, columns[k].value(o))
			}
		}
		rows = rows + row + "\n"
	}

	_, _ = fmt.Fprintln(w, rows)

	return w
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
		if output.APIVersion.isRemovedIn(instance.TargetVersions) {

			removals = removals + 1
		}
		if output.APIVersion.isDeprecatedIn(instance.TargetVersions) {
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
