package finder

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindFiles gets a list of all the files in the directory.
func FindFiles(rootPath string) error {
	var files []string

	root := rootPath
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return err
	}
	for _, file := range files {
		fmt.Println(file)
	}
	return nil
}
