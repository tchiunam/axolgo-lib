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
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"

	"github.com/tchiunam/axolgo-lib/util"
)

// A blockchain transaction
type Transaction struct {
	ID      []byte
	Inputs  []TXInput
	Outputs []TXOutput
}

// A transaction input
type TXInput struct {
	ID  []byte
	Out int
	Sig string
}

// A transaction output
type TXOutput struct {
	Value  int
	PubKey string
}

// Set the ID of the transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	util.PanicOnError(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// Make a transaction to the given address
func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coinbase transaction to %s", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{100, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

func NewTransaction(from string, to string, amount int, chain *BlockChain) (*Transaction, error) {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	if acc < amount {
		return nil, fmt.Errorf("Not enough funds to make a transaction")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		util.PanicOnError(err)

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TXOutput{amount, to})

	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx, nil
}

// Check if the transaction is a Coinbase transaction
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// Check if the transaction input can be unlocked with the given data
func (in *TXInput) CanUnlock(data string) bool {
	return in.Sig == data
}

// Check if the transaction output can be unlocked with the given data
func (out *TXOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
