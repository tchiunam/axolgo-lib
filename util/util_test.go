/*
Copyright © 2022 tchiunam

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package util

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestExpandPath calls ExpandPath with a string,
// checking for a value with tilde expanded.
func TestExpandPath(t *testing.T) {
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	cases := map[string]struct {
		path       string
		expectPath string
	}{
		"nil input": {
			path:       "",
			expectPath: "",
		},
		"only tilde": {
			path:       "~",
			expectPath: homeDir,
		},
		"tilde with path": {
			path:       "~/subdir/sub-subdir",
			expectPath: fmt.Sprintf("%s/subdir/sub-subdir", homeDir),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			path := ExpandPath(c.path)
			assert.Equal(t, c.expectPath, path, "ExpandPath(%v) = %v, want %v", c.path, path, c.expectPath)
		})
	}
}

// TestExtractFileNameWithoutExtension calls ExtractFileNameWithoutExtension
// with a string, checking for a value with the extension removed.
func TestExtractFileNameWithoutExtension(t *testing.T) {
	cases := map[string]struct {
		path           string
		expectFileName string
	}{
		"nil input": {
			path:           "",
			expectFileName: "",
		},
		"no extension": {
			path:           "/etc/path/file",
			expectFileName: "file",
		},
		"extension": {
			path:           "/etc/path/file.txt",
			expectFileName: "file",
		},
		"extension with multiple dots": {
			path:           "/etc/path/file.tar.gz",
			expectFileName: "file",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			fileName := ExtractFileNameWithoutExtension(c.path)
			assert.Equal(
				t,
				c.expectFileName,
				fileName,
				"ExtractFileNameWithoutExtension(%v) = %v, want %v", c.path, fileName, c.expectFileName)
		})
	}
}

// TestUpdateFilePath calls UpdateFilePath with a string,
// checking for a value with customized params.
func TestUpdateFilePath(t *testing.T) {
	cases := map[string]struct {
		path              string
		fileNamePrefix    string
		fileNameSuffix    string
		fileNameExtension string
		expectPath        string
	}{
		"nil input": {
			path:              "",
			fileNamePrefix:    "",
			fileNameSuffix:    "",
			fileNameExtension: "",
			expectPath:        "",
		},
		"no extension": {
			path:              "/etc/path/file",
			fileNamePrefix:    "prefix_",
			fileNameSuffix:    "_suffix",
			fileNameExtension: "",
			expectPath:        "/etc/path/prefix_file_suffix",
		},
		"replace extension and suffix": {
			path:              "/etc/path/file.txt",
			fileNamePrefix:    "",
			fileNameSuffix:    "_suffix",
			fileNameExtension: ".ini",
			expectPath:        "/etc/path/file_suffix.ini",
		},
		"extension with multiple dots and prefix": {
			path:              "/etc/path/file.tar.gz",
			fileNamePrefix:    "prefix_",
			fileNameSuffix:    "",
			fileNameExtension: "",
			expectPath:        "/etc/path/prefix_file.tar.gz",
		},
		"replace extension with multiple dots": {
			path:              "/etc/path/file.tar.gz",
			fileNamePrefix:    "prefix_",
			fileNameSuffix:    "",
			fileNameExtension: ".ini",
			expectPath:        "/etc/path/prefix_file.ini",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			path := UpdateFilePath(
				c.path,
				UpdateFilePathParams{
					FileNamePrefix:    c.fileNamePrefix,
					FileNameSuffix:    c.fileNameSuffix,
					FileNameExtension: c.fileNameExtension,
				},
			)
			assert.Equal(
				t,
				c.expectPath,
				path,
				"UpdateFilePath(%v, %v, %v, %v) = %v, want %v",
				c.path,
				c.fileNamePrefix,
				c.fileNameSuffix,
				c.fileNameExtension,
				path,
				c.expectPath)
		})
	}
}

// TestAddPrefixToFileName calls AddPrefixToFileName with a string,
// checking for a value with the prefix added.
func TestAddPrefixToFileName(t *testing.T) {
	cases := map[string]struct {
		path       string
		prefix     string
		expectPath string
	}{
		"nil input": {
			path:       "",
			prefix:     "",
			expectPath: "",
		},
		"no extension": {
			path:       "/etc/path/file",
			prefix:     "prefix_",
			expectPath: "/etc/path/prefix_file",
		},
		"extension": {
			path:       "/etc/path/file.txt",
			prefix:     "prefix_",
			expectPath: "/etc/path/prefix_file.txt",
		},
		"extension with multiple dots": {
			path:       "/etc/path/file.tar.gz",
			prefix:     "prefix_",
			expectPath: "/etc/path/prefix_file.tar.gz",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			path := AddPrefixToFileName(c.path, c.prefix)
			assert.Equal(
				t,
				c.expectPath,
				path,
				"AddPrefixToFileName(%v, %v) = %v, want %v", c.path, c.prefix, path, c.expectPath)
		})
	}
}

// TestAddSuffixToFileName calls AddSuffixToFileName with a string,
// checking for a value with the suffix added.
func TestAddSuffixToFileName(t *testing.T) {
	cases := map[string]struct {
		path       string
		suffix     string
		expectPath string
	}{
		"nil input": {
			path:       "",
			suffix:     "",
			expectPath: "",
		},
		"no extension": {
			path:       "/etc/path/file",
			suffix:     "_suffix",
			expectPath: "/etc/path/file_suffix",
		},
		"extension": {
			path:       "/etc/path/file.txt",
			suffix:     "_suffix",
			expectPath: "/etc/path/file_suffix.txt",
		},
		"extension with multiple dots": {
			path:       "/etc/path/file.tar.gz",
			suffix:     "_suffix",
			expectPath: "/etc/path/file_suffix.tar.gz",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			path := AddSuffixToFileName(c.path, c.suffix)
			assert.Equal(
				t,
				c.expectPath,
				path,
				"AddSuffixToFileName(%v, %v) = %v, want %v", c.path, c.suffix, path, c.expectPath)
		})
	}
}

// TestIntToHex calls IntToHex with an int64 value, checking for a
// string with the hex value.
func TestIntToHex(t *testing.T) {
	cases := map[string]struct {
		num int64
	}{
		"normal input": {
			num: 100,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := IntToHex(c.num)
			assert.NoError(t, err, "IntToHex(%v) = %v, want %v", c.num, err, nil)
		})
	}
}
