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

// TestGeneratePassphrase calls GeneratePassphrase and checks the
// generated passphrase.
func TestGeneratePassphrase(t *testing.T) {
	cases := map[string]struct {
		length       int
		expectLength int
	}{
		"length == 10": {
			length:       10,
			expectLength: 10 * 2,
		},
		"length == 50": {
			length:       50,
			expectLength: 50 * 2,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := GeneratePassphrase(c.length)
			assert.Nil(t, err, "GeneratePassphrase(%v) = %v, want nil", c.length, err)
			assert.Equal(t, c.expectLength, len(actual), "GeneratePassphrase(%v) = %v, want %v", c.expectLength, actual, c.expectLength)
		})
	}
}

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

// MockWithCryptographyOptionsError is a mock implementation of CryptographyOptions
// that can be used for testing error.
func MockWithCryptographyOptionsError(v string) CryptographyOptionsFunc {
	return func(o *CryptographyOptions) error {
		o.OutputFilename = v
		return fmt.Errorf("mock error")
	}
}

// TestEncryptDecrypt tests the Encrypt and Decrypt functions
func TestEncryptDecrypt(t *testing.T) {
	cases := map[string]struct {
		data       []byte
		passphrase string
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

			// Test with cryptography options
			_, err = Encrypt(c.data, c.passphrase, WithHashFunc(CreateHash))
			assert.Nil(t, err, "Encrypt(%x, %v) = %v, want nil", c.data, c.passphrase, err)
		})
	}
}

// TestEncryptDecryptInvalid tests the Encrypt and Decrypt functions with invalid input
func TestEncryptDecryptInvalid(t *testing.T) {
	cases := map[string]struct {
		optFns     func(*CryptographyOptions) error
		data       []byte
		passphrase string
	}{
		"normal input": {
			data:       []byte("The quick brown fox jumps over the lazy dog"),
			passphrase: "iamthebest",
			optFns:     MockWithCryptographyOptionsError("foo.txt"),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Encrypt(c.data, c.passphrase, c.optFns)
			assert.Error(t, err, "Encrypt(%x, %v) = %v, want error", c.data, c.passphrase, err)
			_, err = Decrypt(c.data, c.passphrase, c.optFns)
			assert.Error(t, err, "Decrypt(%x, %v) = %v, want error", c.data, c.passphrase, err)
		})
	}
}

// TestDecryptPanic calls Decrypt and checks the panic
func TestDecryptPanic(t *testing.T) {
	cases := map[string]struct {
		optFns     func(*CryptographyOptions) error
		data       []byte
		passphrase string
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

func MockHashFuncWithError(input string) string {
	myarr := [1]string{"foo"}
	index := 2

	// This is a hack to get around the fact that the compiler doesn't know that
	// the array is of length 1.
	return myarr[index]
}

func TestEncryptDecryptWithHashInvalid(t *testing.T) {
	cases := map[string]struct {
		data       []byte
		passphrase string
		hashFunc   PassphraseHashFunc
	}{
		"nil data": {
			data:       []byte("The quick brown fox jumps over the lazy dog"),
			passphrase: "iamthebest",
			hashFunc:   MockHashFuncWithError,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Panics(
				t,
				func() { Encrypt(c.data, c.passphrase, WithHashFunc(c.hashFunc)) },
				"Decrypt(%x, %v, %v), want panic, %v", c.data, c.passphrase, c.hashFunc)
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
		_cleanTestEncryptDecryptFile(c.encOutputFilename, c.decOutputFilename)
	}
}

func TestEncryptDecryptFileInvalid(t *testing.T) {
	cases := map[string]struct {
		filename             string
		encFilename          string
		encOutputFilename    string
		decOutputFilename    string
		passphrase           string
		optFn                func(string) CryptographyOptionsFunc
		expectEncErrorString string
		expectDecErrorString string
	}{
		"invalid option": {
			filename:             filepath.Join("testdata", "story.txt"),
			encOutputFilename:    filepath.Join("testdata", "foo.txt"),
			decOutputFilename:    filepath.Join("testdata", "foo.txt"),
			passphrase:           "iamthebest",
			optFn:                MockWithCryptographyOptionsError,
			expectEncErrorString: "Fail to read cryptography options: mock error",
			expectDecErrorString: "Fail to read cryptography options: mock error",
		},
		"input file not exists": {
			filename:             filepath.Join("testdata", "bar.txt"),
			encOutputFilename:    filepath.Join("testdata", "foo.txt"),
			decOutputFilename:    filepath.Join("testdata", "story.txt.dec"),
			passphrase:           "iamthebest",
			optFn:                WithOutputFilename,
			expectEncErrorString: "open testdata/bar.txt: no such file or directory",
			expectDecErrorString: "open testdata/foo.txt: no such file or directory",
		},
		"invalid encOutputFilename": {
			filename:             filepath.Join("testdata", "story.txt"),
			encOutputFilename:    filepath.Join("testdata", ""),
			decOutputFilename:    filepath.Join("testdata", "story.txt.dec"),
			passphrase:           "iamthebest",
			optFn:                WithOutputFilename,
			expectEncErrorString: "open testdata: is a directory",
			expectDecErrorString: "read testdata: is a directory",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := EncryptFile(c.filename, c.passphrase, c.optFn(c.encOutputFilename))
			assert.Error(t, err, "EncryptFile(%v, %v, %v) = %v, want error", c.filename, c.passphrase, c.encOutputFilename, err)
			assert.Equal(
				t,
				c.expectEncErrorString,
				err.Error(),
				"EncryptFile(%v, %v, %v) = %v, want %v", c.filename, c.passphrase, c.encOutputFilename, err, c.expectEncErrorString)
			_, err = DecryptFile(c.encOutputFilename, c.passphrase, c.optFn(c.decOutputFilename))
			assert.Error(t, err, "DecryptFile(%v, %v, %v) = %v, want error", c.encOutputFilename, c.passphrase, c.decOutputFilename, err)
			assert.Equal(
				t,
				c.expectDecErrorString,
				err.Error(),
				"DecryptFile(%v, %v, %v) = %v, want %v", c.encOutputFilename, c.passphrase, c.decOutputFilename, err, c.expectDecErrorString)
		})
		_cleanTestEncryptDecryptFile(c.encOutputFilename, c.decOutputFilename)
	}
}
