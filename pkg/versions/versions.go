package verisions

// APIVersion is a file that has an apiVersion in it
type APIVersion struct {
	// Name is the name of the api version
	Name string
	// Deprecated is a boolean that indicates if the apiVersion is deprecated
	Deprecated bool
}

// APIVersionList is a map of apiVersions and if they are deprecated or not.
// TODO: Find a better way to generate this.
var APIVersionList = []APIVersion{
	{"extensions/v1beta1", true},
	{"apps/v1", false},
}
