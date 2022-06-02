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
