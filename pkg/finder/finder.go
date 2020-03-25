// Copyright 2020 Fairwinds
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package finder

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"k8s.io/klog"

	"github.com/fairwindsops/api-version-finder/pkg/api"
)

// Dir is the finder dirlication
type Dir struct {
	RootPath string
	FileList []string
	APIFiles []*File
}

// File is a file that has an apiVersion in it
type File struct {
	Name       string       `json:"file,omitempty" yaml:"file,omitempty"`
	APIVersion *api.Version `json:"api,omitempty" yaml:"api,omitempty"`
}

// NewFinder returns a new struct with config portions complete.
func NewFinder(path string) *Dir {
	cfg := &Dir{}
	if path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			klog.Fatal(err)
		}
		cfg.RootPath = cwd
	} else {
		cfg.RootPath = path
	}
	return cfg
}

// FindVersions runs the finder. This will populate the
// dir struct with any files that might be versioned.
func (dir *Dir) FindVersions() error {
	err := dir.listFiles()
	if err != nil {
		return err
	}
	err = dir.scanFiles()
	if err != nil {
		return err
	}
	if dir.APIFiles != nil {
		for _, file := range dir.APIFiles {
			klog.V(6).Infof("%t - %s - %s", file.APIVersion.Deprecated, file.APIVersion.Name, file.Name)
		}
	}
	return nil
}

// listFiles gets a list of all the files in the directory.
func (dir *Dir) listFiles() error {
	var files []string

	if _, err := os.Stat(dir.RootPath); os.IsNotExist(err) {
		return fmt.Errorf("specified path does not exist")
	}
	err := filepath.Walk(dir.RootPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	dir.FileList = files
	return nil
}

// scanFiles loops through the file list and finds versioned files
// to add to the dir struct
func (dir *Dir) scanFiles() error {

	for _, file := range dir.FileList {
		klog.V(8).Infof("processing file: %s", file)
		apiFile, err := checkForAPIVersion(file)
		if err != nil {
			klog.V(2).Infof("error scanning file %s: %s", file, err.Error())
		}
		if apiFile != nil {
			dir.APIFiles = append(dir.APIFiles, apiFile)
		}
	}
	return nil
}

// checkForAPIVersion checks a filename to see if
// it is an api-versioned Kubernetes object.
// Returns the File object if it is.
func checkForAPIVersion(file string) (*File, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	version, err := api.IsVersioned(data)
	if err != nil {
		return nil, err
	}
	if version == nil {
		return nil, nil
	}
	versionedFile := &File{
		Name:       file,
		APIVersion: version,
	}
	return versionedFile, nil
}
