package util

import (
	"testing"
)

func TestSplitFilename_Success(t *testing.T) {
	tests := []struct {
		path         string
		expectedName string
		expectedExt  string
	}{
		{"/foo/bar/somefile.log", "somefile", "log"},
		{"/path/to/file.txt", "file", "txt"},
		{"document.pdf", "document", "pdf"},
		{"./test.json", "test", "json"},
	}

	for _, tt := range tests {
		name, ext, err := SplitFilename(tt.path)
		if err != nil {
			t.Errorf("SplitFilename(%q) unexpected error: %v", tt.path, err)
			continue
		}
		if name != tt.expectedName {
			t.Errorf("SplitFilename(%q) name = %q, want %q", tt.path, name, tt.expectedName)
		}
		if ext != tt.expectedExt {
			t.Errorf("SplitFilename(%q) ext = %q, want %q", tt.path, ext, tt.expectedExt)
		}
	}
}

func TestSplitFilename_Errors(t *testing.T) {
	tests := []struct {
		path        string
		expectedErr string
	}{
		{"", "filename cannot be empty"},
		{"/", "filename cannot be empty"},
		{"/path/to/", "file extension cannot be empty"}, // "to" is valid name but no extension
		{"filename", "file extension cannot be empty"},
		{"/path/to/filename", "file extension cannot be empty"},
		{".hidden", "filename cannot be empty"}, // .hidden files have empty name
	}

	for _, tt := range tests {
		name, ext, err := SplitFilename(tt.path)
		if err == nil {
			t.Errorf("SplitFilename(%q) expected error %q, got name=%q ext=%q", tt.path, tt.expectedErr, name, ext)
			continue
		}
		if err.Error() != tt.expectedErr {
			t.Errorf("SplitFilename(%q) error = %q, want %q", tt.path, err.Error(), tt.expectedErr)
		}
	}
}
