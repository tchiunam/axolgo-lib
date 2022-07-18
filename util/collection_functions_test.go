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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Predicate function type for testing predicate function
// with string array.
type TestStringPredicateFunction1 func(string) bool
type TestStringPredicateFunction2 func(string) string

// TestAny4s calls Any4s to check for valid return values.
func TestAny4s(t *testing.T) {
	cases := map[string]struct {
		input             []string
		predicateFunction TestStringPredicateFunction1
		expectBool        bool
	}{
		"empty strings without matching": {
			input:             []string{},
			predicateFunction: func(s string) bool { return len(s) > 3 },
			expectBool:        false,
		},
		"at least one matches": {
			input:             []string{"dog", "cat", "mouse", "bird", "fish"},
			predicateFunction: func(s string) bool { return len(s) > 3 },
			expectBool:        true,
		},
		"no matches": {
			input:             []string{"car", "bus", "truck", "train", "boat"},
			predicateFunction: func(s string) bool { return len(s) > 5 },
			expectBool:        false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := Any4s(c.input, c.predicateFunction)
			assert.Equal(t, c.expectBool, actual, "Any4s(%v, %v) = %v, want %v", c.input, c.predicateFunction, actual, c.expectBool)
		})
	}
}

// TestAll4s calls All4s to check for valid return values.
func TestAll4s(t *testing.T) {
	cases := map[string]struct {
		input             []string
		predicateFunction TestStringPredicateFunction1
		expectBool        bool
	}{
		"empty strings without matching": {
			input:             []string{},
			predicateFunction: func(s string) bool { return len(s) > 3 },
			expectBool:        false,
		},
		"all match": {
			input:             []string{"dog", "cat", "mouse", "bird", "fish"},
			predicateFunction: func(s string) bool { return len(s) > 2 },
			expectBool:        true,
		},
		"at least one not match": {
			input:             []string{"car", "bus", "truck", "train", "boat"},
			predicateFunction: func(s string) bool { return len(s) > 4 },
			expectBool:        false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := All4s(c.input, c.predicateFunction)
			assert.Equal(t, c.expectBool, actual, "All4s(%v, %v) = %v, want %v", c.input, c.predicateFunction, actual, c.expectBool)
		})
	}
}

// TestFilter4s calls Filter4s to check for valid return values.
func TestFilter4s(t *testing.T) {
	cases := map[string]struct {
		input             []string
		predicateFunction TestStringPredicateFunction1
		expectStringArray []string
	}{
		"empty strings without matching": {
			input:             []string{},
			predicateFunction: func(s string) bool { return len(s) > 3 },
			expectStringArray: []string{},
		},
		"at least one matches": {
			input:             []string{"dog", "cat", "mouse", "bird", "fish"},
			predicateFunction: func(s string) bool { return len(s) > 3 },
			expectStringArray: []string{"mouse", "bird", "fish"},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := Filter4s(c.input, c.predicateFunction)
			if !reflect.DeepEqual(actual, c.expectStringArray) {
				t.Errorf("Filter4s(%v, %v) = %v, want %v", c.input, c.predicateFunction, actual, c.expectStringArray)
			}
		})
	}
}

// TestMap4s calls Map4s to check for valid return values.
func TestMap4s(t *testing.T) {
	cases := map[string]struct {
		input             []string
		predicateFunction TestStringPredicateFunction2
		expectStrings     []string
	}{
		"append to string": {
			input:             []string{"dog", "cat", "mouse", "bird", "fish"},
			predicateFunction: func(s string) string { return s + "!" },
			expectStrings:     []string{"dog!", "cat!", "mouse!", "bird!", "fish!"},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := Map4s(c.input, c.predicateFunction)
			if !reflect.DeepEqual(actual, c.expectStrings) {
				t.Errorf("Map4s(%v, %v) = %v, want %v", c.input, c.predicateFunction, actual, c.expectStrings)
			}
		})
	}
}
