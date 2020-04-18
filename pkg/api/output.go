package api

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"gopkg.in/yaml.v2"
)

var padChar = byte(' ')

// DisplayOutput prints the output based on desired variables
func DisplayOutput(outputs []*Output, outputFormat string, showNonDeprecated bool, targetVersion string) error {
	if len(outputs) == 0 {
		fmt.Println("There were no apiVersions found that match our records.")
		return nil
	}
	usableOutputs := filterNonDeprecated(outputs, targetVersion, showNonDeprecated)
	var err error
	var outData []byte
	switch outputFormat {
	case "tabular":
		t, err := tabOut(usableOutputs, targetVersion)
		if err != nil {
			return err
		}
		err = t.Flush()
		if err != nil {
			return err
		}
		return nil
	case "json":
		outData, err = json.Marshal(usableOutputs)
		if err != nil {
			return err
		}
		fmt.Println(string(outData))
	case "yaml":
		outData, err = yaml.Marshal(usableOutputs)
		if err != nil {
			return err
		}
		fmt.Println(string(outData))
	default:
		fmt.Println("output format should be one of (json,yaml,tabular)")
	}
	return nil
}

func filterNonDeprecated(outputs []*Output, targetVersion string, showNonDeprecated bool) []*Output {
	var usableOutputs []*Output

	if showNonDeprecated {
		usableOutputs = outputs
	} else {

		for _, output := range outputs {
			if output.APIVersion.IsDeprecatedIn(targetVersion) {
				usableOutputs = append(usableOutputs, output)
			}
		}
	}
	return usableOutputs
}

func tabOut(outputs []*Output, targetVersion string) (*tabwriter.Writer, error) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 15, 2, padChar, 0)

	if len(outputs) == 0 {
		_, err := fmt.Fprintln(w, "APIVersions were found, but none were deprecated. Try --show-all.")
		if err != nil {
			return nil, err
		}
		return w, nil
	}
	_, err := fmt.Fprintln(w, "KIND\t VERSION\t DEPRECATED\t DEPRECATED IN\t RESOURCE NAME\t")
	if err != nil {
		return nil, err
	}
	for _, output := range outputs {
		kind := output.APIVersion.Kind
		deprecated := fmt.Sprintf("%t", output.APIVersion.IsDeprecatedIn(targetVersion))
		version := output.APIVersion.Name
		fileName := output.Name
		deprecatedIn := output.APIVersion.DeprecatedIn
		if deprecatedIn == "" {
			deprecatedIn = "n/a"
		}

		_, err = fmt.Fprintf(w, "%s\t %s\t %s\t %s\t %s\t\n", kind, version, deprecated, deprecatedIn, fileName)
		if err != nil {
			return nil, err
		}
	}
	return w, nil
}

// GetReturnCode checks for deprecated versions and returns a code.
// takes a boolean to ignore any errors.
func GetReturnCode(outputs []*Output, ignoreErrors bool, targetVersion string) int {
	if ignoreErrors {
		return 0
	}
	for _, output := range outputs {
		if output.APIVersion.IsDeprecatedIn(targetVersion) {
			return 2
		}
	}
	return 0
}
