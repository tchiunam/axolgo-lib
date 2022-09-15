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
	"crypto/rsa"
	"hash"
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

// TestEncryptRSA calls the EncryptRSA function to make sure
// the RSA encryption works.
func TestEncryptRSA(t *testing.T) {
	cases := map[string]struct {
		data      []byte
		publicKey rsa.PublicKey
		hashFunc  hash.Hash
	}{
		"normal input with 2048 bit": {
			data: []byte("hello world"),
			publicKey: func() rsa.PublicKey {
				_, publicKey, _ := GenerateRSAKeyPair(2048)
				return *publicKey
			}(),
			hashFunc: nil,
		},
		"normal input with 4096 bit": {
			data: []byte("hello world"),
			publicKey: func() rsa.PublicKey {
				_, publicKey, _ := GenerateRSAKeyPair(4096)
				return *publicKey
			}(),
			hashFunc: nil,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := EncryptRSA(c.data, c.publicKey, WithOAEPHashFunc(c.hashFunc))
			assert.NoError(t, err, "EncryptRSA(%v, %v) = %v", c.data, c.publicKey, err)
		})
	}
}

// TestEncryptRSAInvalid calls the EncryptRSA function to make sure
// errors are returned when invalid parameters are passed.
func TestEncryptRSAInvalid(t *testing.T) {
	cases := map[string]struct {
		data              []byte
		publicKey         rsa.PublicKey
		optFn             func(*CryptographyOptions) error
		expectErrorString string
	}{
		"errornous hash function": {
			data: []byte("hello world"),
			publicKey: func() rsa.PublicKey {
				_, publicKey, _ := GenerateRSAKeyPair(2048)
				return *publicKey
			}(),
			optFn:             MockWithCryptographyOptionsError("foo"),
			expectErrorString: "Fail to read cryptography options: mock error",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := EncryptRSA(c.data, c.publicKey, c.optFn)
			assert.Error(t, err, "EncryptRSA(%v, %v) = %v", c.data, c.publicKey, err)
			assert.Equal(t, c.expectErrorString, err.Error(), "EncryptRSA(%v, %v) = %v", c.data, c.publicKey, err)
		})
	}
}
