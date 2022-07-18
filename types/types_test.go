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
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParameter calls Parameter to check for initialization
func TestParameter(t *testing.T) {
	paramName := "hammer"
	paramValue := "nail"

	cases := map[string]struct {
		name  *string
		value *string
	}{
		"normal input": {
			// initialize the name and value string pointers
			name:  &paramName,
			value: &paramValue,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// initialize the parameter
			param := Parameter{Name: tc.name, Value: tc.value}
			// check the name and value
			assert.Equal(t, *tc.name, *param.Name, "Expected name %s, got %s", *tc.name, *param.Name)
			assert.Equal(t, *tc.value, *param.Value, "Expected value %s, got %s", *tc.value, *param.Value)
		})
	}
}

// TestStringArrayFlag calls StringArrayFlag to check for initialization
func TestStringArrayFlag(t *testing.T) {
	cases := map[string]struct {
		value []string
	}{
		"two values": {
			value: []string{"lion", "tiger"},
		},
		"four values": {
			value: []string{"lion", "tiger", "bear", "wolf"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// loop through the values and append them to the flag
			flag := StringArrayFlag{}
			for _, value := range tc.value {
				flag.Set(value)
			}
			// loop through the value string array
			for i, v := range flag {
				assert.Equal(t, tc.value[i], v, "Expected value %s, got %s", tc.value[i], v)
			}
			// check the flag String() result
			assert.Equal(
				t,
				fmt.Sprintf("%v", tc.value),
				flag.String(),
				"Expected value %s, got %s", fmt.Sprintf("%v", tc.value), flag.String())
		})
	}
}
