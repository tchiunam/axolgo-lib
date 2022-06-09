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
	"fmt"
	"strings"
)

type AxolVarTemplateStringForAttribute string

// Enum values for AxolVarTemplateStringForAttribute
const (
	AxolVarTemplateStringFieldWithKeyValue AxolVarTemplateStringForAttribute = "{{range $i, $v := .VAR0}}{{if $i}}, {{end}}{{`{`}}{{.VAR1}}: {{.VAR2}}{{`}`}}{{end}}"
)

// Generate a new template string using the given values
func (a AxolVarTemplateStringForAttribute) New(values ...string) string {
	var oldNewStrings []string
	for i, v := range values {
		oldNewStrings = append(oldNewStrings, fmt.Sprintf("VAR%d", i), v)
	}
	replacer := strings.NewReplacer(oldNewStrings...)
	return replacer.Replace(string(a))
}
