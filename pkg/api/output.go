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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/olekukonko/tablewriter"

	"gopkg.in/yaml.v3"
)

var padChar = byte(' ')

// Output is a thing that has an apiVersion in it
type Output struct {
	// Name is the name of the object in question.
	// This might be an object name, or a release
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// FilePath is the full path of the file if the output came from a file
	FilePath string `json:"filePath,omitempty" yaml:"filePath,omitempty"`
	// Namespace is the namespace that the object is in
	// The output may resolve this to UNKNOWN if there is no way of determining it
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	// APIVersion is the version object corresponding to this output
	APIVersion *Version `json:"api,omitempty" yaml:"api,omitempty"`
	// Deprecated is a boolean indicating whether or not the version is deprecated
	Deprecated bool `json:"deprecated" yaml:"deprecated"`
	// Removed is a boolean indicating whether or not the version has been removed
	Removed bool `json:"removed" yaml:"removed"`
	// ReplacementAvailable is a boolean indicating whether or not the replacement is available
	ReplacementAvailable bool `json:"replacementAvailable" yaml:"replacementAvailable"`
	// CustomColumns is a list of column headers to be displayed with -ocustom or -omarkdown
	CustomColumns []string `json:"-" yaml:"-"`
}

// Instance is an instance of the API. This holds configuration for a "run" of Pluto
type Instance struct {
	Outputs                       []*Output         `json:"items,omitempty" yaml:"items,omitempty"`
	IgnoreDeprecations            bool              `json:"-" yaml:"-"`
	IgnoreRemovals                bool              `json:"-" yaml:"-"`
	IgnoreUnavailableReplacements bool              `json:"-" yaml:"-"`
	OnlyShowRemoved               bool              `json:"-" yaml:"-"`
	NoHeaders                     bool              `json:"-" yaml:"-"`
	OutputFormat                  string            `json:"-" yaml:"-"`
	TargetVersions                map[string]string `json:"target-versions,omitempty" yaml:"target-versions,omitempty"`
	DeprecatedVersions            []Version         `json:"-" yaml:"-"`
	CustomColumns                 []string          `json:"-" yaml:"-"`
	Components                    []string          `json:"-" yaml:"-"`
}

// DisplayOutput prints the output based on desired variables
func (instance *Instance) DisplayOutput() error {
	if len(instance.Outputs) == 0 && (instance.OutputFormat == "normal" || instance.OutputFormat == "wide") {
		fmt.Println("There were no resources found with known deprecated apiVersions.")
		return nil
	}

	instance.FilterOutput()
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
	case "custom":
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
	case "markdown":
		var c columnList
		if len(instance.CustomColumns) >= 1 {
			c = instance.customColumns()
		} else {
			c = instance.wideColumns()
		}
		t := instance.markdownOut(c)
		if t != nil {
			t.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			t.SetCenterSeparator("|")
			t.Render()
		}
	case "csv":
		var c columnList
		if len(instance.CustomColumns) >= 1 {
			c = instance.customColumns()
		} else {
			c = instance.wideColumns()
		}
		csvWriter, err := instance.csvOut(c)
		if err != nil {
			return err
		}

		csvWriter.Flush()

		err = csvWriter.Error()
		if err != nil {
			return err
		}
	}
	return nil
}

// FilterOutput filters the outputs that get printed
// first it fills out the Deprecated and Removed booleans
// then it returns the outputs that are either deprecated or removed
// and in the component list
// additionally, if instance.OnlyShowDeprecated is true, it will remove the
// apiVersions that are deprecated but not removed
func (instance *Instance) FilterOutput() {
	var usableOutputs []*Output
	for _, output := range instance.Outputs {
		output.Deprecated = output.APIVersion.isDeprecatedIn(instance.TargetVersions)
		output.Removed = output.APIVersion.isRemovedIn(instance.TargetVersions)
		output.ReplacementAvailable = output.APIVersion.isReplacementAvailableIn(instance.TargetVersions)
		switch instance.OnlyShowRemoved {
		case false:
			if output.Deprecated || output.Removed {
				if StringInSlice(output.APIVersion.Component, instance.Components) {
					usableOutputs = append(usableOutputs, output)
				}
			}
		case true:
			if output.Removed {
				if StringInSlice(output.APIVersion.Component, instance.Components) {
					usableOutputs = append(usableOutputs, output)
				}
			}
		}
	}
	instance.Outputs = usableOutputs
}

// removeDeprecatedOnly is a list replacement operation
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

	if !instance.NoHeaders {
		var headers string
		for _, k := range columnIndexes {
			if k == 0 {
				headers = fmt.Sprintf("%s\t", columns[k].header())
			} else {
				headers = fmt.Sprintf("%s %s\t", headers, columns[k].header())
			}
		}
		_, _ = fmt.Fprintln(w, headers)
	}

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

func (instance *Instance) markdownOut(columns columnList) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)

	if len(instance.Outputs) == 0 {
		_, _ = fmt.Println("No output to display")
		return nil
	}

	columnIndexes := make([]int, 0, len(columns))
	for k := range columns {
		columnIndexes = append(columnIndexes, k)
	}
	sort.Ints(columnIndexes)

	if !instance.NoHeaders {
		var headers []string
		for _, k := range columnIndexes {
			headers = append(headers, columns[k].header())
		}

		table.SetHeader(headers)
	}

	for _, o := range instance.Outputs {
		var row []string
		for _, k := range columnIndexes {
			row = append(row, columns[k].value(o))
		}
		table.Append(row)
	}

	return table
}

func (instance *Instance) csvOut(columns columnList) (*csv.Writer, error) {
	csvWriter := csv.NewWriter(os.Stdout)

	if len(instance.Outputs) == 0 {
		_, _ = fmt.Println("No output to display")
	}

	columnIndexes := make([]int, 0, len(columns))
	for k := range columns {
		columnIndexes = append(columnIndexes, k)
	}
	sort.Ints(columnIndexes)

	var csvData [][]string

	if !instance.NoHeaders {
		var headers []string
		for _, k := range columnIndexes {
			headers = append(headers, columns[k].header())
		}

		csvData = append(csvData, headers)
	}

	for _, o := range instance.Outputs {
		var row []string
		for _, k := range columnIndexes {
			row = append(row, columns[k].value(o))
		}
		csvData = append(csvData, row)
	}

	for i := range csvData {
		if err := csvWriter.Write(csvData[i]); err != nil {
			return nil, err
		}
	}

	return csvWriter, nil
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
			if output.APIVersion.isReplacementAvailableIn(instance.TargetVersions) || !instance.IgnoreUnavailableReplacements {
				deprecations = deprecations + 1
			}
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
