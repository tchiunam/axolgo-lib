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
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// PassphraseHashFunc is a function that returns a hash of a passphrase.
type PassphraseHashFunc func(string) string

// CryptographyOptionsFunc is a type alias for CryptographyOptions functional option
type CryptographyOptionsFunc func(*CryptographyOptions) error

// CryptographyOptions are discrete set of options that are valid for loading the
// configuration that is used to encrypt/decrypt files.
type CryptographyOptions struct {
	CustomHashFunc PassphraseHashFunc
	OutputFilename string
}

// WithCustomHashFunc is a helper function to construct functional options
// that sets a custom hash function for the passphrase.
func WithCustomHashFunc(fn PassphraseHashFunc) CryptographyOptionsFunc {
	return func(o *CryptographyOptions) error {
		o.CustomHashFunc = fn
		return nil
	}
}

// WithOutputFilename is a helper function to construct functional options
// that sets the output filename for the encrypted/decrypted file.
func WithOutputFilename(v string) CryptographyOptionsFunc {
	return func(o *CryptographyOptions) error {
		o.OutputFilename = v
		return nil
	}
}

// Generate a random passphrase. Returns a random passphrase and an error if any.
func GeneratePassphrase(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Create a hash from a string
func CreateHash(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Evaluate the functional options and set the options in the CryptographyOptions struct
func evaluateCryptographyInputOptions(options *CryptographyOptions, optFns ...CryptographyOptionsFunc) error {
	for _, optFn := range optFns {
		if err := optFn(options); err != nil {
			return fmt.Errorf("Fail to read cryptography options: %v", err)
		}
	}

	return nil
}

// Encrypt data with a passphrase. Nonce is created by the function.
// Returns the encrypted data and an error if any.
func Encrypt(data []byte, passphrase string, optFns ...CryptographyOptionsFunc) ([]byte, error) {
	options := CryptographyOptions{CustomHashFunc: CreateHash}
	if err := evaluateCryptographyInputOptions(&options, optFns...); err != nil {
		return nil, err
	}

	block, _ := aes.NewCipher([]byte(options.CustomHashFunc(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt data with a passphrase. Returns the decrypted
// data and an error if any.
func Decrypt(data []byte, passphrase string, optFns ...CryptographyOptionsFunc) ([]byte, error) {
	options := CryptographyOptions{CustomHashFunc: CreateHash}
	if err := evaluateCryptographyInputOptions(&options, optFns...); err != nil {
		return nil, err
	}

	key := []byte(options.CustomHashFunc(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// Encrypt a file with a passphrase. Returns an error if any.
// filename is the path to the file to be encrypted.
// passphrase is the passphrase to use to encrypt the file.
func EncryptFile(filename string, passphrase string, optFns ...CryptographyOptionsFunc) ([]byte, error) {
	return _cryptFile(Encrypt, filename, passphrase, optFns...)
}

// Decrypt a file with a passphrase. Returns an error if any.
// filename is the path to the file to be decrypted.
// passphrase is the passphrase to use to decrypt the file.
func DecryptFile(filename string, passphrase string, optFns ...CryptographyOptionsFunc) ([]byte, error) {
	return _cryptFile(Decrypt, filename, passphrase, optFns...)
}

// fn is the function to be used to crypt/decrypt the file.
func _cryptFile(
	fn func([]byte, string, ...CryptographyOptionsFunc) ([]byte, error),
	filename string,
	passphrase string,
	optFns ...CryptographyOptionsFunc) ([]byte, error) {
	var options CryptographyOptions
	if err := evaluateCryptographyInputOptions(&options, optFns...); err != nil {
		return nil, err
	}

	if content, err := os.ReadFile(filename); err == nil {
		if data, err := fn(content, passphrase, optFns...); err == nil {
			if options.OutputFilename != "" {
				if err = os.WriteFile(options.OutputFilename, data, 0644); err != nil {
					return nil, err
				}
			}
			return data, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
