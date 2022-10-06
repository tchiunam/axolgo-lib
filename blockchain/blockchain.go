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
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"os"

	"github.com/dgraph-io/badger"
	"github.com/tchiunam/axolgo-lib/util"
)

// Data in the genesis block
const genesisData = "First block in the chain - Genesis"

// BlockChain structure
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

// An iterator for iterating the blockchain in database
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// Check if blockchain database exists
func DBExists(dbPath string) bool {
	if _, err := os.Stat(dbPath + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}

	return true
}

// InitBlockChain creates a new blockchain with a genesis block
func InitBlockChain(dbPath string, address string) *BlockChain {
	var lastHash []byte

	if DBExists(dbPath) {
		panic("Blockchain already exists.")
	}

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	util.PanicOnError(err)

	err = db.Update(func(txn *badger.Txn) error {
		abTx := CoinbaseTx(address, genesisData)
		genesis := Genesis(abTx)
		err = txn.Set(genesis.Hash, genesis.Serialize())
		util.PanicOnError(err)

		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash

		return err
	})

	util.PanicOnError(err)

	chain := BlockChain{lastHash, db}
	return &chain
}

// Continue a blockchain by pulling the last hash
func ContinueBlockChain(dbPath string) *BlockChain {
	if DBExists(dbPath) == false {
		panic("No existing blockchain found. Create one first.")
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	util.PanicOnError(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		util.PanicOnError(err)

		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})

		return err
	})

	util.PanicOnError(err)

	chain := BlockChain{lastHash, db}
	return &chain
}

// AddBlock is a helper function that adds a new block to the chain using
// the previous block's hash
func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		util.PanicOnError(err)

		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})

		return err
	})

	util.PanicOnError(err)

	newBlock := CreateBlock(transactions, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		util.PanicOnError(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
}

// Iterator returns a BlockChainIterator that can be used to iterate over
// the blockchain.
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}
	return iter
}

// Next returns the next block in the blockchain
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		util.PanicOnError(err)

		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})

		return err
	})

	util.PanicOnError(err)

	iter.CurrentHash = block.PrevHash

	return block
}

// Find all transactions for a given address that is unspent
func (chain *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTxs []Transaction
	spentTxOutputs := make(map[string][]int)
	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTxOutputs[txID] != nil {
					for _, spentOut := range spentTxOutputs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.ID)
						spentTxOutputs[inTxID] = append(spentTxOutputs[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTxs
}

// Find all unspent transaction outputs for a given address
func (chain *BlockChain) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	unspentTxs := chain.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTxs {
		for _, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// Find transaction outputs for a given address that can be
// used as inputs for a new transaction
func (chain *BlockChain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

// FindTransaction finds a transaction by its ID
func (chain *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

// SignTransaction signs inputs of a Transaction
func (chain *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTxs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTx, err := chain.FindTransaction(in.ID)
		util.PanicOnError(err)
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}

	tx.Sign(privKey, prevTxs)
}
