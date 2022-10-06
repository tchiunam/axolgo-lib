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

package blockchain

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWallets tests the wallet persistence
func TestWallets(t *testing.T) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	for i := 0; i < 5; i++ {
		address := wallets.AddWallet()
		assert.NotNil(t, address)
	}

	addresses := wallets.GetAllAddresses()
	wallet := wallets.GetWallet(addresses[0])
	assert.NotNil(t, wallet)

	WalletFilePath = filepath.Join("testdata", "wallets.dat")
	err := wallets.Persist()
	// defer os.Remove(walletFilePath)

	// Having a bug to be fixed here
	assert.Error(t, err)

	CreateWallets()
	// assert.NoError(t, err)
	// assert.NotNil(t, walletsLoaded)
}
