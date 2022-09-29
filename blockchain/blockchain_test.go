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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Clean up the test database
func _cleanTestBadgerDatabase(
	dbPath string) {
	// Delete the database files from the dbPath
	if _, err := os.Stat(dbPath); err == nil {
		os.RemoveAll(dbPath)
	}
}

// TestBlockChain tests transaction operations on the blockchain
func TestBlockChain(t *testing.T) {
	cases := map[string]struct {
		from   string
		to     string
		amount int
	}{
		"tx 1": {
			from:   "John",
			to:     "Jane",
			amount: 20,
		},
		"tx 2": {
			from:   "John",
			to:     "Mary",
			amount: 30,
		},
		"tx 3": {
			from:   "Mary",
			to:     "Nancy",
			amount: 10,
		},
		"tx 4": {
			from:   "Nancy",
			to:     "John",
			amount: 5,
		},
	}

	dbPath := filepath.Join("testdata", "db", "blockchain")
	os.MkdirAll(dbPath, 0755)
	defer _cleanTestBadgerDatabase(dbPath)

	chain := InitBlockChain(dbPath, "John")
	defer chain.Database.Close()
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			tx, _ := NewTransaction(c.from, c.to, c.amount, chain)
			assert.NotPanics(t, func() { chain.AddBlock([]*Transaction{tx}) }, "AddBlock should not panic")
		})
	}
	// Close the database connection so that we can open it again
	chain.Database.Close()

	t.Run("Verify John's balance", func(t *testing.T) {
		chain := ContinueBlockChain(dbPath)

		balance := 0
		UTXOs := chain.FindUTXO("John")

		for _, out := range UTXOs {
			balance += out.Value
		}

		assert.Equal(t, 55, balance, "Balance of John should be 55")
		// Close the database connection so that we can open it again
		chain.Database.Close()
	})

	t.Run("Verify PoW", func(t *testing.T) {
		chain := ContinueBlockChain(dbPath)

		iterator := chain.Iterator()
		for {
			block := iterator.Next()
			pow := NewProof(block)
			assert.True(t, pow.Validate(), "Proof of work is not valid")

			if len(block.PrevHash) == 0 {
				break
			}
		}
	})
}
