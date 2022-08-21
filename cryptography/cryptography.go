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
	"io/ioutil"
)

// CryptographyOptionsFunc is a type alias for CryptographyOptions functional option
type CryptographyOptionsFunc func(*CryptographyOptions) error

// CryptographyOptions are discrete set of options that are valid for loading the
// configuration that is used to encrypt/decrypt files.
type CryptographyOptions struct {
	OutputFilename string
}

// WithOutputFilename is a helper function to construct functional options
// that sets the output filename for the encrypted/decrypted file.
func WithOutputFilename(v string) CryptographyOptionsFunc {
	return func(o *CryptographyOptions) error {
		o.OutputFilename = v
		return nil
	}
}

// Create a hash from a string
func CreateHash(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Encrypt data with a passphrase. Nonce is created by the function.
// Returns the encrypted data and an error if any.
func Encrypt(data []byte, passphrase string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(CreateHash(passphrase)))
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
func Decrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(CreateHash(passphrase))
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
		panic(err.Error())
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
	fn func([]byte, string) ([]byte, error),
	filename string,
	passphrase string,
	optFns ...CryptographyOptionsFunc) ([]byte, error) {
	var options CryptographyOptions
	for _, optFn := range optFns {
		if err := optFn(&options); err != nil {
			return nil, fmt.Errorf("Fail to read cryptography options: %v", err)
		}
	}

	if content, err := ioutil.ReadFile(filename); err == nil {
		if data, err := fn(content, passphrase); err == nil {
			if options.OutputFilename != "" {
				if err = ioutil.WriteFile(options.OutputFilename, data, 0644); err != nil {
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
