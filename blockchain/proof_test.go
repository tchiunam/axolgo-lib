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

// TestProofOfWork tests the proof of work algorithm
func TestProofOfWork(t *testing.T) {
	cases := map[string]struct {
		data []string
	}{
		"normal input": {
			data: []string{"First block", "Second block", "Third block"},
		},
	}

	dbPath := filepath.Join("testdata", "db", "proof-of-work")
	os.MkdirAll(dbPath, 0755)
	defer _cleanTestBadgerDatabase(dbPath)

	chain := InitBlockChain(dbPath)
	defer chain.Database.Close()
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			for _, d := range c.data {
				chain.AddBlock(d)
			}
		})
	}

	t.Run("Iterating all blocks", func(t *testing.T) {
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
