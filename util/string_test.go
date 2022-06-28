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

import "testing"

// TestHushedStringPtr calls HushedStringPtr with a string,
// checking for a valid return value.
func TestHushedStringPtr(t *testing.T) {
	testString := "something"

	cases := map[string]struct {
		StringPtr    *string
		ExpectString string
	}{
		"nil input": {
			StringPtr:    nil,
			ExpectString: "",
		},
		"non-nil input": {
			StringPtr:    &testString,
			ExpectString: "something",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := HushedStringPtr(c.StringPtr)
			if *actual != c.ExpectString {
				t.Errorf("HushedStringPtr(%v) = %v, want %v", c.StringPtr, *actual, c.ExpectString)
			}
		})
	}
}
