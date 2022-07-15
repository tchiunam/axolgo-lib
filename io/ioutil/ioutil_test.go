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

package ioutil

import (
	"strings"
	"testing"
)

// TestReadIniFile calls ReadIniFile to check for reading
// of INI file.
func TestReadIniFile(t *testing.T) {
	cases := map[string]struct {
		filepath           string
		expectSectionNames []string
	}{
		"valid file": {
			filepath:           "./testdata/ioutil_test.ini",
			expectSectionNames: []string{"book", "airplane"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			file, err := ReadIniFile(tc.filepath)
			if err != nil {
				t.Errorf("Expected no error, got %s", err.Error())
			}
			for _, sectionName := range tc.expectSectionNames {
				if section := file.Section(sectionName); section.Name() != sectionName {
					t.Errorf("Expected section name %s, got %s", sectionName, section.Name())
				}
			}
		})
	}
}

// TestReadIniFileInvalid calls ReadIniFile to check for error cases
func TestReadIniFileInvalid(t *testing.T) {
	cases := map[string]struct {
		filepath            string
		expectStringInError string
	}{
		"non-exist file": {
			filepath:            "./non-exist-file.ini",
			expectStringInError: "open ./non-exist-file.ini: no such file or directory",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := ReadIniFile(tc.filepath)
			if err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.expectStringInError) {
				t.Errorf("Expected error containing %s, got %s", tc.expectStringInError, err.Error())
			}
		})
	}
}

// Structure of a book author
type TestReadYamlFileInvalidBookAuthor struct {
	Name string `mapstructure:"name"`
}

// Structure of book attributes
type TestReadYamlFileInvalidBookAttributes struct {
	Title  string                            `mapstructure:"title"`
	Author TestReadYamlFileInvalidBookAuthor `mapstructure:"author"`
}

// Structure of a book
type TestReadYamlFileInvalidInventory struct {
	Book TestReadYamlFileInvalidBookAttributes `mapstructure:"book"`
}

// TestReadYamlFile calls ReadYamlFile to check
// for reading of YAML file.
func TestReadYamlFile(t *testing.T) {
	testBook := TestReadYamlFileInvalidInventory{
		Book: TestReadYamlFileInvalidBookAttributes{
			Title: "The Hitchhiker's Guide to the Galaxy",
			Author: TestReadYamlFileInvalidBookAuthor{
				Name: "Douglas Adams",
			},
		},
	}

	cases := map[string]struct {
		filepath string
		expect   TestReadYamlFileInvalidInventory
	}{
		"valid file": {
			filepath: "./testdata/ioutil_test.yaml",
			expect:   testBook,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var inventory TestReadYamlFileInvalidInventory
			options := WithCFOClass(&inventory)
			_, err := ReadYamlFile(tc.filepath, options)

			if err != nil {
				t.Errorf("Expected no error, got %s", err.Error())
			}
			if inventory.Book.Title != tc.expect.Book.Title {
				t.Errorf("Expected title %s, got %s", tc.expect.Book.Title, inventory.Book.Title)
			}
			if inventory.Book.Author.Name != tc.expect.Book.Author.Name {
				t.Errorf("Expected author name %s, got %s", tc.expect.Book.Author.Name, inventory.Book.Author.Name)
			}
		})
	}
}

// TestReadYamlFileInvalid calls ReadYamlFile to check for error cases
func TestReadYamlFileInvalid(t *testing.T) {
	cases := map[string]struct {
		filepath            string
		expectStringInError string
	}{
		"file does not exist": {
			filepath:            "./non-exist-file.ini",
			expectStringInError: "open ./non-exist-file.ini: no such file or directory",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := ReadYamlFile(tc.filepath)
			if err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.expectStringInError) {
				t.Errorf("Expected error containing %s, got %s", tc.expectStringInError, err.Error())
			}
		})
	}
}
