package util

import (
	"errors"
	"path/filepath"
	"strings"
)

// SplitFilename splits a path into the filename and extension.
// For example, "/path/to/file.txt" becomes ("file", "txt").
// Returns an error if either the name or extension are empty.
func SplitFilename(path string) (name, ext string, err error) {
	// Get just the filename from the full path
	filename := filepath.Base(path)

	// Check for invalid filenames first
	if filename == "" || filename == "." || filename == ".." || filename == "/" {
		return "", "", errors.New("filename cannot be empty")
	}

	// Get the extension (includes the dot)
	ext = filepath.Ext(filename)

	// Remove the extension from filename to get the name
	name = strings.TrimSuffix(filename, ext)

	// Remove the dot from extension
	ext = strings.TrimPrefix(ext, ".")

	// Check for empty name or extension
	if name == "" {
		return "", "", errors.New("filename cannot be empty")
	}
	if ext == "" {
		return "", "", errors.New("file extension cannot be empty")
	}

	return name, ext, nil
}
