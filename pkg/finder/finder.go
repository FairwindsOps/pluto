package finder

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/klog"
)

// App is the finder application
type App struct {
	RootPath    *string
	FileList    []*string
	APIVersions []*APIVersion
}

// APIVersion is a file that has an apiVersion in it
type APIVersion struct {
	// FileName is the name of the file that contains the apiVersion
	FileName string
	// Name is the name of the api version
	Name string
	// Type is the type of file that contains the version
	Type *string
	// Deprecated is a boolean that indicates if the apiVersion is deprecated
	Deprecated bool
}

// APIVersionList is a map of apiVersions and if they are deprecated or not.
// TODO: Find a better way to generate this.
var APIVersionList = map[string]bool{
	"extensions/v1beta1": true,
	"apps/v1":            false,
}

// NewFinder returns a new struct with config portions complete.
func NewFinder(path string) *App {
	cfg := &App{}
	if path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			klog.Fatal(err)
		}
		cfg.RootPath = &cwd
	} else {
		cfg.RootPath = &path
	}
	return cfg
}

// Execute runs the finder
func (app *App) Execute() error {
	err := app.findFiles()
	if err != nil {
		return err
	}
	err = app.scanFiles()
	if err != nil {
		return err
	}
	if app.APIVersions != nil {
		for _, version := range app.APIVersions {
			klog.Infof("%t - %s - %s", version.Deprecated, version.Name, version.FileName)
		}
	}
	return nil
}

// FindFiles gets a list of all the files in the directory.
func (app *App) findFiles() error {
	var files []*string

	err := filepath.Walk(*app.RootPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, &path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	app.FileList = files
	return nil
}

func (app *App) scanFiles() error {
	if app.FileList == nil {
		return fmt.Errorf("files list is nil")
	}

	for _, file := range app.FileList {
		klog.V(8).Infof("processing file: %s", *file)
		version, err := checkForAPIVersion(*file)
		if err != nil {
			klog.Errorf("error scanning file %s: %s", file, err.Error())
		}
		if version != nil {
			app.APIVersions = append(app.APIVersions, version)
		}
	}
	return nil
}

func checkForAPIVersion(file string) (*APIVersion, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		version, deprecated := lineContainsVersion(scanner.Text())
		if version != nil {
			api := &APIVersion{
				FileName:   file,
				Name:       *version,
				Deprecated: deprecated,
			}
			return api, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func lineContainsVersion(line string) (*string, bool) {
	for version, deprecated := range APIVersionList {
		if strings.Contains(line, version) {
			return &version, deprecated
		}
	}
	return nil, false
}
