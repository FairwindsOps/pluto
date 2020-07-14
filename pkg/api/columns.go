package api

import "fmt"

// Column is an interface for printing columns
type column interface {
	header() string
	value(output *Output) string
}

type columnList map[int]column

//PossibleColumnNames is the list of implmented columns
var PossibleColumnNames = []string{
	"NAME",
	"NAMESPACE",
	"KIND",
	"VERSION",
	"REPLACEMENT",
	"DEPRECATED",
	"DEPRECATED IN",
	"REMOVED",
	"REMOVED IN",
}

var possibleColumns = []column{
	new(name),
	new(namespace),
	new(kind),
	new(version),
	new(replacement),
	new(deprecated),
	new(deprecatedIn),
	new(removed),
	new(removedIn),
}

// name is the output name
type name struct{}

func (n name) header() string              { return "NAME" }
func (n name) value(output *Output) string { return output.Name }

// namespace is the output namespace if available
type namespace struct{}

func (ns namespace) header() string { return "NAMESPACE" }
func (ns namespace) value(output *Output) string {
	if output.Namespace == "" {
		return "<UNKNOWN>"
	}
	return output.Namespace
}

// kind is the output apiVersion kind
type kind struct{}

func (k kind) header() string              { return "KIND" }
func (k kind) value(output *Output) string { return output.APIVersion.Kind }

// version is the output apiVersion
type version struct{}

func (v version) header() string              { return "VERSION" }
func (v version) value(output *Output) string { return output.APIVersion.Name }

// replacement is the output replacement apiVersion
type replacement struct{}

func (r replacement) header() string              { return "REPLACEMENT" }
func (r replacement) value(output *Output) string { return output.APIVersion.ReplacementAPI }

// deprecated is the output for the boolean Deprecated
type deprecated struct{}

func (d deprecated) header() string              { return "DEPRECATED" }
func (d deprecated) value(output *Output) string { return fmt.Sprintf("%t", output.Deprecated) }

// removed is the output for the boolean Deprecated
type removed struct{}

func (r removed) header() string              { return "REMOVED" }
func (r removed) value(output *Output) string { return fmt.Sprintf("%t", output.Removed) }

// deprecatedIn is the string value of when an output was deprecated
type deprecatedIn struct{}

func (di deprecatedIn) header() string              { return "DEPRECATED IN" }
func (di deprecatedIn) value(output *Output) string { return output.APIVersion.DeprecatedIn }

// removedIn is the string value of when an output was deprecated
type removedIn struct{}

func (ri removedIn) header() string              { return "REMOVED IN" }
func (ri removedIn) value(output *Output) string { return output.APIVersion.RemovedIn }

// normalColumns returns the list of columns for -onormal
func (instance *Instance) normalColumns() columnList {
	columnList := columnList{
		0: new(name),
		1: new(kind),
		2: new(version),
		3: new(replacement),
		4: new(removed),
		5: new(deprecated),
	}
	return columnList
}

// wideColumns returns the list of columns for -owide
func (instance *Instance) wideColumns() columnList {
	columnList := columnList{
		0: new(name),
		1: new(namespace),
		2: new(kind),
		3: new(version),
		4: new(replacement),
		5: new(deprecated),
		6: new(deprecatedIn),
		7: new(removed),
		8: new(removedIn),
	}
	return columnList
}

// customColumns returns a custom list of columns based on names
func (instance *Instance) customColumns() columnList {
	var outputColumns = make(map[int]column)
	for _, d := range instance.CustomColumns {
		for i, c := range possibleColumns {
			if d == c.header() {
				outputColumns[i] = c
			}
		}
	}
	return outputColumns
}
