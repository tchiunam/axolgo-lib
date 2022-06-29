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
)

// TestExpandPath calls ExpandPath with a string,
// checking for a value with tilde expanded.
func TestExpandPath(t *testing.T) {
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	cases := map[string]struct {
		Path       string
		ExpectPath string
	}{
		"nil input": {
			Path:       "",
			ExpectPath: "",
		},
		"only tilde": {
			Path:       "~",
			ExpectPath: homeDir,
		},
		"tilde with path": {
			Path:       "~/subdir/sub-subdir",
			ExpectPath: fmt.Sprintf("%s/subdir/sub-subdir", homeDir),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			path := ExpandPath(c.Path)
			if path != c.ExpectPath {
				t.Errorf("ExpandPath(%q) = %q, want %q, case %q", c.Path, path, c.ExpectPath, name)
			}
		})
	}
}
