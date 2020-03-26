package api

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"gopkg.in/yaml.v2"
)

// DisplayOutput prints the output based on desired variables
func DisplayOutput(outputs []*Output, outputFormat string, showNonDeprecated bool) error {
	if len(outputs) == 0 {
		fmt.Println("There were no apiVersions found that match our records.")
		return nil
	}
	var err error
	var outData []byte
	var usableOutputs []*Output
	switch outputFormat {
	case "tabular":
		if showNonDeprecated {
			usableOutputs = outputs
		} else {
			for _, output := range outputs {
				if output.APIVersion.Deprecated {
					usableOutputs = append(usableOutputs, output)
				}
			}
		}
		if len(usableOutputs) == 0 {
			fmt.Println("APIVersions were found, but none were deprecated. Try --show-non-deprecated.")
			return nil
		}
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 2, ' ', 0)
		_, err = fmt.Fprintln(w, "KIND\t VERSION\t DEPRECATED\t RESOURCE NAME")
		if err != nil {
			return err
		}
		for _, output := range usableOutputs {
			kind := output.APIVersion.Kind
			deprecated := fmt.Sprintf("%t", output.APIVersion.Deprecated)
			version := output.APIVersion.Name
			fileName := output.Name

			_, err = fmt.Fprintf(w, "%s\t %s\t %s\t %s\t\n", kind, version, deprecated, fileName)
			if err != nil {
				return err
			}
		}
		err = w.Flush()
		if err != nil {
			return err
		}
	case "json":
		outData, err = json.Marshal(outputs)
		if err != nil {
			return err
		}
	case "yaml":
		outData, err = yaml.Marshal(outputs)
		if err != nil {
			return err
		}
	default:
		fmt.Println("output format should be one of (json,yaml,tabular)")
	}

	fmt.Println(string(outData))
	return nil
}
