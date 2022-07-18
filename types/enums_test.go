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

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAxolVarTemplateStringFieldWithKeyValue calls the New function
// of AxolVarTemplateStringFieldWithKeyValue, checking for a valid
// return value.
func TestAxolVarTemplateStringFieldWithKeyValueNew(t *testing.T) {
	cases := map[string]struct {
		input        []string
		expectString string
	}{
		"non-nil input": {
			input:        []string{"0RAV", "1RAV", "2RAV"},
			expectString: "{{range $i, $v := .0RAV}}{{if $i}}, {{end}}{{`{`}}{{.1RAV}}: {{.2RAV}}{{`}`}}{{end}}",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			a := AxolVarTemplateStringFieldWithKeyValue.New(c.input...)
			assert.Equal(t, c.expectString, a, "Expected %s, got %s", c.expectString, a)
		})
	}
}
