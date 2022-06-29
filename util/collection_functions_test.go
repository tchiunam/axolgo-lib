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
)

// Predicate function type for testing predicate function
// with string array.
type TestStringPredicateFunction func(string) bool

// TestAny4s calls Any4s to check for valid return values.
func TestAny4s(t *testing.T) {
	cases := map[string]struct {
		Strings           []string
		PredicateFunction TestStringPredicateFunction
		ExpectBool        bool
	}{
		"empty strings without matching": {
			Strings:           []string{},
			PredicateFunction: func(s string) bool { return len(s) > 3 },
			ExpectBool:        false,
		},
		"at least one matches": {
			Strings:           []string{"dog", "cat", "mouse", "bird", "fish"},
			PredicateFunction: func(s string) bool { return len(s) > 3 },
			ExpectBool:        true,
		},
		"no matches": {
			Strings:           []string{"car", "bus", "truck", "train", "boat"},
			PredicateFunction: func(s string) bool { return len(s) > 5 },
			ExpectBool:        false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := Any4s(c.Strings, c.PredicateFunction)
			if actual != c.ExpectBool {
				t.Errorf("Any4s(%v, %v) = %v, want %v, case %q", c.Strings, c.PredicateFunction, actual, c.ExpectBool, name)
			}
		})
	}
}

// TestAll4s calls All4s to check for valid return values.
func TestAll4s(t *testing.T) {
	cases := map[string]struct {
		Strings           []string
		PredicateFunction TestStringPredicateFunction
		ExpectBool        bool
	}{
		"empty strings without matching": {
			Strings:           []string{},
			PredicateFunction: func(s string) bool { return len(s) > 3 },
			ExpectBool:        false,
		},
		"all match": {
			Strings:           []string{"dog", "cat", "mouse", "bird", "fish"},
			PredicateFunction: func(s string) bool { return len(s) > 2 },
			ExpectBool:        true,
		},
		"at least one not match": {
			Strings:           []string{"car", "bus", "truck", "train", "boat"},
			PredicateFunction: func(s string) bool { return len(s) > 4 },
			ExpectBool:        false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := All4s(c.Strings, c.PredicateFunction)
			if actual != c.ExpectBool {
				t.Errorf("All4s(%v, %v) = %v, want %v, case %q", c.Strings, c.PredicateFunction, actual, c.ExpectBool, name)
			}
		})
	}
}

// TestFilter4s calls Filter4s to check for valid return values.
func TestFilter4s(t *testing.T) {
	cases := map[string]struct {
		Strings           []string
		PredicateFunction TestStringPredicateFunction
		ExpectStringArray []string
	}{
		"empty strings without matching": {
			Strings:           []string{},
			PredicateFunction: func(s string) bool { return len(s) > 3 },
			ExpectStringArray: []string{},
		},
		"at least one matches": {
			Strings:           []string{"dog", "cat", "mouse", "bird", "fish"},
			PredicateFunction: func(s string) bool { return len(s) > 3 },
			ExpectStringArray: []string{"mouse", "bird", "fish"},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := Filter4s(c.Strings, c.PredicateFunction)
			if !reflect.DeepEqual(actual, c.ExpectStringArray) {
				t.Errorf("Filter4s(%v, %v) = %v, want %v, case %q", c.Strings, c.PredicateFunction, actual, c.ExpectStringArray, name)
			}
		})
	}
}
