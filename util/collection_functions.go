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

// Returns true if any of the items in the given slice
// satisfies the predicate f.
func Any4s(ss []string, f func(string) bool) bool {
	for _, v := range ss {
		if f(v) {
			return true
		}
	}
	return false
}

// Returns true if all strings in the slice satisfy the
// given predicate f.
func All4s(ss []string, f func(string) bool) bool {
	for _, v := range ss {
		if !f(v) {
			return false
		}
	}
	return true
}

// Returns a new slice of strings containing all strings
// in the given slice that satisfy the given predicate f.
func Filter4s(ss []string, f func(string) bool) []string {
	nss := make([]string, 0)
	for _, v := range ss {
		if f(v) {
			nss = append(nss, v)
		}
	}
	return nss
}

// Map collection function for string slice.
// Function f is applied to all strings in the given slice
// and a new slice containing the results is returned.
func Map4s(ss []string, f func(string) string) []string {
	nss := make([]string, len(ss))
	for i, v := range ss {
		nss[i] = f(v)
	}
	return nss
}
