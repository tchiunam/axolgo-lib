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
	"bytes"
	"encoding/binary"
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/mr-tron/base58"
)

// Expand tilde to home directory
func ExpandPath(path string) string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	if path == "~" {
		return homeDir
	} else if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, (path)[2:])
	}

	return path
}

// Extract the file name without the extension
func ExtractFileNameWithoutExtension(path string) string {
	if path == "" {
		return ""
	}

	// Remove the base path in case the file name has no extension
	// which causes it not to be processed in the next step.
	path = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	for strings.Contains(path, ".") {
		path = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}
	return path
}

// Parameters for the `UpdateFilePath` function
type UpdateFilePathParams struct {
	FileNamePrefix    string
	FileNameSuffix    string
	FileNameExtension string // If empty, the original extension will be used
}

// Update the file path with file name prefix, file name suffix
// or other customizations.
func UpdateFilePath(path string, params UpdateFilePathParams) string {
	if path == "" {
		return ""
	}

	fileName := filepath.Base(path)
	f := strings.Split(fileName, ".")
	newFileName := params.FileNamePrefix
	newFileName += f[0]
	newFileName += params.FileNameSuffix

	if params.FileNameExtension == "" {
		// Append the extension (ex. .tar.gz) to the file name
		for i := 1; i < len(f); i++ {
			newFileName = fmt.Sprintf("%s.%s", newFileName, f[i])
		}
	} else {
		newFileName = fmt.Sprintf("%s%s", newFileName, params.FileNameExtension)
	}

	return filepath.Join(filepath.Dir(path), newFileName)
}

// Add prefix to file name
func AddPrefixToFileName(path string, prefix string) string {
	return UpdateFilePath(path, UpdateFilePathParams{FileNamePrefix: prefix})
}

// Add suffix to file name
func AddSuffixToFileName(path string, suffix string) string {
	return UpdateFilePath(path, UpdateFilePathParams{FileNameSuffix: suffix})
}

// Covert int to bytes of hex
func IntToHex(num int64) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// Panic if there is error
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// Base58 encoding
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

// Base58 decoding
func Base58Decode(input []byte) ([]byte, error) {
	decode, err := base58.Decode(string(input[:]))
	if err != nil {
		return nil, err
	}

	return decode, nil
}
