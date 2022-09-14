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

package cryptography

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateRSAKeyPair calls the GenerateRSAKeyPair function
// to make sure it returns a key pair.
func TestGenerateRSAKeyPair(t *testing.T) {
	cases := map[string]struct {
		bits int
	}{
		"2048 bits": {
			bits: 2048,
		},
		"4096 bits": {
			bits: 2048,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			privateKey, publicKey, err := GenerateRSAKeyPair(c.bits)
			assert.NoError(t, err, "GenerateRSAKeyPair(%v) = %v", c.bits, err)
			assert.NotNil(t, privateKey, "private key should not be nil")
			assert.NotNil(t, publicKey, "public key should not be nil")
		})
	}
}

// TestGenerateRSAKeyPairInvalid calls the GenerateRSAKeyPair function
// to make sure errors are returned when invalid parameters are passed.
func TestGenerateRSAKeyPairInvalid(t *testing.T) {
	cases := map[string]struct {
		bits int
	}{
		"1 bits": {
			bits: 1,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, _, err := GenerateRSAKeyPair(c.bits)
			assert.Error(t, err, "GenerateRSAKeyPair(%v) = %v", c.bits, err)
		})
	}
}
