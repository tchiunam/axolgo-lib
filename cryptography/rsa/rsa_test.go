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

package rsa

import (
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tchiunam/axolgo-lib/cryptography"
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

// TestEncryptDecryptRSA calls the EncryptRSA and DecryptRSA function
// to make sure they return the original data.
func TestEncryptDecryptRSA(t *testing.T) {
	cases := map[string]struct {
		data       []byte
		privateKey *rsa.PrivateKey
	}{
		"normal input with 2048 bit": {
			data: []byte("hello world"),
			privateKey: func() *rsa.PrivateKey {
				privateKey, _, _ := GenerateRSAKeyPair(2048)
				return privateKey
			}(),
		},
		"normal input with 4096 bit": {
			data: []byte("hello world"),
			privateKey: func() *rsa.PrivateKey {
				privateKey, _, _ := GenerateRSAKeyPair(4096)
				return privateKey
			}(),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			encryptedBytes, err := EncryptRSA(
				c.data,
				c.privateKey.PublicKey,
				cryptography.WithOAEPHashFunc(sha256.New()))
			assert.NoError(t, err, "EncryptRSA(%v, %v) = %v", string(c.data), c.privateKey.PublicKey, err)
			decryptedBytes, err := DecryptRSA(encryptedBytes, c.privateKey)
			assert.NoError(t, err, "DecryptRSA(%v, %v) = %v", encryptedBytes, c.privateKey, err)
			assert.Equal(
				t,
				string(c.data),
				string(decryptedBytes),
				"%v and %v should be equal", string(c.data), string(decryptedBytes))
		})
	}
}

// MockWithCryptographyOptionsError is a mock implementation of CryptographyOptions
// that can be used for testing error.
func MockWithCryptographyOptionsError(v string) cryptography.CryptographyOptionsFunc {
	return func(o *cryptography.CryptographyOptions) error {
		o.OutputFilename = v
		return fmt.Errorf("mock error")
	}
}

// TestEncryptDecryptRSAInvalid calls the EncryptRSA and DecryptRSA
// function to make sure errors are returned when invalid parameters are passed.
func TestEncryptDecryptRSAInvalid(t *testing.T) {
	cases := map[string]struct {
		data                 []byte
		privateKey           *rsa.PrivateKey
		optFn                func(*cryptography.CryptographyOptions) error
		expectEncErrorString string
		expectDecErrorString string
	}{
		"errornous hash function": {
			data: []byte("hello world"),
			privateKey: func() *rsa.PrivateKey {
				privateKey, _, _ := GenerateRSAKeyPair(2048)
				return privateKey
			}(),
			optFn:                MockWithCryptographyOptionsError("foo"),
			expectEncErrorString: "Fail to read cryptography options: mock error",
			expectDecErrorString: "Fail to read cryptography options: mock error",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := EncryptRSA(c.data, c.privateKey.PublicKey, c.optFn)
			assert.Error(t, err, "EncryptRSA(%v, %v) = %v", string(c.data), c.privateKey.PublicKey, err)
			assert.Equal(
				t,
				c.expectEncErrorString,
				err.Error(),
				"EncryptRSA(%v, %v) = %v", string(c.data), c.privateKey.PublicKey, err)
			// Doesn't matter what the encrypted data is, as long as it's not nil
			_, err = DecryptRSA(c.data, c.privateKey, c.optFn)
			assert.Error(t, err, "DecryptRSA(%v, %v) = %v", string(c.data), c.privateKey, err)
			assert.Equal(
				t,
				c.expectDecErrorString,
				err.Error(),
				"DecryptRSA(%v, %v) = %v", string(c.data), c.privateKey, err)
		})
	}
}

// TestSignRSA calls the SignRSA to make sure it can sign data
func TestSignRSA(t *testing.T) {
	cases := map[string]struct {
		data       []byte
		privateKey *rsa.PrivateKey
	}{
		"normal input with 2048 bit": {
			data: []byte("hello world"),
			privateKey: func() *rsa.PrivateKey {
				privateKey, _, _ := GenerateRSAKeyPair(2048)
				return privateKey
			}(),
		},
		"normal input with 4096 bit": {
			data: []byte("hello world"),
			privateKey: func() *rsa.PrivateKey {
				privateKey, _, _ := GenerateRSAKeyPair(4096)
				return privateKey
			}(),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := SignRSA(c.data, c.privateKey)
			assert.NoError(t, err, "SignRSA(%v, %v) = %v", c.data, c.privateKey, err)
		})
	}
}

// TestSignRSAInvalid calls the SignRSA to make sure errors are returned
func TestSignRSAInvalid(t *testing.T) {
	cases := map[string]struct {
		data              []byte
		privateKey        *rsa.PrivateKey
		optFn             func(*cryptography.CryptographyOptions) error
		expectErrorString string
	}{
		"errornous hash function": {
			data: []byte("hello world"),
			privateKey: func() *rsa.PrivateKey {
				privateKey, _, _ := GenerateRSAKeyPair(2048)
				return privateKey
			}(),
			optFn:             MockWithCryptographyOptionsError("foo"),
			expectErrorString: "Fail to read cryptography options: mock error",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := SignRSA(c.data, c.privateKey, c.optFn)
			assert.Error(t, err, "SignRSA(%v, %v) = %v", string(c.data), c.privateKey, err)
			assert.Equal(
				t,
				c.expectErrorString,
				err.Error(),
				"SignRSA(%v, %v) = %v", string(c.data), c.privateKey, err)
		})
	}
}
