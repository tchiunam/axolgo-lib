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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateHash calls CreateHash and checks the output
func TestCreateHash(t *testing.T) {
	cases := map[string]struct {
		input        string
		expectString string
	}{
		"string1": {
			input:        "^M@lGg*N9AcAiKv8R7$*Iv*7D",
			expectString: "0239e5621ee962d96d8aed6735df1a4d",
		},
		"string2": {
			input:        "@69O^7hp7iTR1Nj#vJ94#4Tphiy!P&TAxDY",
			expectString: "f0852b81dbde6d1119fff905b9ee6714",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := CreateHash(c.input)
			assert.Equal(t, c.expectString, actual, "CreateHash(%v) = %v, want %v", c.input, actual, c.expectString)
		})
	}
}

// TestEncryptDecrypt tests the Encrypt and Decrypt functions
func TestEncryptDecrypt(t *testing.T) {
	cases := map[string]struct {
		data         []byte
		passphrase   string
		expectString string
	}{
		"normal input": {
			data:       []byte("The quick brown fox jumps over the lazy dog"),
			passphrase: "iamthebest",
		},
		"normal input 2": {
			data:       []byte("Have you ever seen a caterpillar eat an apple?"),
			passphrase: "notthebest",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actualEncrypt, err := Encrypt(c.data, c.passphrase)
			assert.Nil(t, err, "Encrypt(%x, %v) = %v, want nil", c.data, c.passphrase, err)
			actualDecrypt, err := Decrypt(actualEncrypt, c.passphrase)
			assert.Nil(t, err, "Decrypt(%x, %v) = %v, want nil", actualEncrypt, c.passphrase, err)
			assert.Equal(t, c.data, actualDecrypt, "Decrypt(%v, %v) = %v, want %v", actualEncrypt, c.passphrase, actualDecrypt, c.data)
		})
	}
}

// TestEncryptDecryptInvalid call Decrypt with an invalid data
func TestDecryptInvalid(t *testing.T) {
	cases := map[string]struct {
		data         []byte
		passphrase   string
		expectString string
	}{
		"nil data": {
			data:       nil,
			passphrase: "iamthebest",
		},
		"empty data": {
			data:       []byte(""),
			passphrase: "notthebest",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Panics(t, func() { Decrypt(c.data, c.passphrase) }, "Decrypt(%x, %v), want panic, %v", c.data, c.passphrase)
		})
	}
}

// Clean up the test files
func _cleanTestEncryptDecryptFile(
	encOutputFilename string,
	decOutputFilename string) {
	// Delete the encrypted file if it exists
	if _, err := os.Stat(encOutputFilename); err == nil {
		os.Remove(encOutputFilename)
	}
	// Delete the decrypted file if it exists
	if _, err := os.Stat(decOutputFilename); err == nil {
		os.Remove(decOutputFilename)
	}
}

// TestEncryptDecryptFile tests the EncryptFile and DecryptFile functions
func TestEncryptDecryptFile(t *testing.T) {
	cases := map[string]struct {
		filename          string
		encFilename       string
		encOutputFilename string
		decOutputFilename string
		passphrase        string
		optFns            func(*CryptographyOptions) error
	}{
		"normal input": {
			filename:          filepath.Join("testdata", "story.txt"),
			encOutputFilename: filepath.Join("testdata", "story.txt.enc"),
			decOutputFilename: filepath.Join("testdata", "story.txt.dec"),
			passphrase:        "iamthebest",
		},
	}

	for name, c := range cases {
		_cleanTestEncryptDecryptFile(c.encOutputFilename, c.decOutputFilename)
		t.Run(name, func(t *testing.T) {
			_, err := EncryptFile(c.filename, c.passphrase, WithOutputFilename(c.encOutputFilename))
			assert.Nil(t, err, "EncryptFile(%v, %v, %v) = %v, want nil", c.filename, c.passphrase, c.encOutputFilename, err)
			_, err = DecryptFile(c.encOutputFilename, c.passphrase, WithOutputFilename(c.decOutputFilename))
			assert.Nil(t, err, "DecryptFile(%v, %v, %v) = %v, want nil", c.encFilename, c.passphrase, c.decOutputFilename, err)
		})
	}
}

// MockWithCryptographyOptionsError is a mock implementation of CryptographyOptions
// that can be used for testing error.
func MockWithCryptographyOptionsError(v string) CryptographyOptionsFunc {
	return func(o *CryptographyOptions) error {
		o.OutputFilename = v
		return fmt.Errorf("mock error")
	}
}

func TestEncryptDecryptFileInvalid(t *testing.T) {
	cases := map[string]struct {
		filename          string
		encFilename       string
		encOutputFilename string
		decOutputFilename string
		passphrase        string
		optFns            func(*CryptographyOptions) error
	}{
		"invalid option": {
			optFns:            MockWithCryptographyOptionsError("foo.txt"),
			filename:          filepath.Join("testdata", "story.txt"),
			encOutputFilename: filepath.Join("testdata", "story.txt.enc"),
			decOutputFilename: filepath.Join("testdata", "story.txt.dec"),
			passphrase:        "iamthebest",
		},
		"file not exists": {
			optFns:            WithOutputFilename("foo.txt"),
			filename:          filepath.Join("testdata", "bar.txt"),
			encOutputFilename: filepath.Join("testdata", "story.txt.enc"),
			decOutputFilename: filepath.Join("testdata", "story.txt.dec"),
			passphrase:        "iamthebest",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := EncryptFile(c.filename, c.passphrase, c.optFns)
			assert.Error(t, err, "EncryptFile(%v, %v) = %v, want error", c.filename, c.passphrase, err)
			_, err = DecryptFile(c.encOutputFilename, c.passphrase, c.optFns)
			assert.Error(t, err, "DecryptFile(%v, %v) = %v, want error", c.encOutputFilename, c.passphrase, err)
		})
		_cleanTestEncryptDecryptFile(c.encOutputFilename, c.decOutputFilename)
	}
}
