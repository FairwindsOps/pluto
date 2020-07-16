package api

import (
	"os"
)

// IsFileOrStdin detects if a file exists, or returns true if - is passed
func IsFileOrStdin(name string) bool {
	if name == "-" {
		return true
	}
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// StringInSlice returns true if the string is contained in the slice
func StringInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if s == v {
			return true
		}
	}
	return false
}
