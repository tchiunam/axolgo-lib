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
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// Generate a new RSA key pair
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

// Encrypt a message using RSA public key
func EncryptRSA(data []byte, publicKey rsa.PublicKey, optFns ...CryptographyOptionsFunc) ([]byte, error) {
	options := CryptographyOptions{OAEPHashFunc: sha256.New()}
	if err := options.merge(optFns...); err != nil {
		return nil, err
	}

	if encryptedBytes, err := rsa.EncryptOAEP(
		options.OAEPHashFunc,
		rand.Reader,
		&publicKey,
		data,
		nil); err == nil {
		return encryptedBytes, nil
	} else {
		return nil, err
	}
}

// Decrypt a message using RSA public key
func DecryptRSA(data []byte, privateKey *rsa.PrivateKey, optFns ...CryptographyOptionsFunc) ([]byte, error) {
	options := CryptographyOptions{OAEPHashFunc: sha256.New()}
	if err := options.merge(optFns...); err != nil {
		return nil, err
	}

	if decryptedBytes, err := rsa.DecryptOAEP(
		options.OAEPHashFunc,
		rand.Reader,
		privateKey,
		data,
		nil); err == nil {
		return decryptedBytes, nil
	} else {
		return nil, err
	}
}

// Sign a message using RSA private key
func SignRSA(data []byte, privateKey *rsa.PrivateKey, optFns ...CryptographyOptionsFunc) ([]byte, error) {
	options := CryptographyOptions{OAEPHashFunc: sha256.New()}
	if err := options.merge(optFns...); err != nil {
		return nil, err
	}

	msgHash := options.OAEPHashFunc
	_, err := msgHash.Write(data)
	if err != nil {
		return nil, err
	}
	msgHashSum := msgHash.Sum(nil)

	return rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
}
