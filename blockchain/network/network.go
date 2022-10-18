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

package network

import (
	"fmt"
	"os"
	"syscall"

	"github.com/vrecan/death/v3"

	"github.com/tchiunam/axolgo-lib/blockchain"
)

const (
	protocol        = "tcp"
	protocolVersion = 1  // Bump whenever a backwards-incompatible protocol change is made
	commandLength   = 12 // The length of all commands in bytes
)

var (
	nodeAddress     string
	mineAddress     string
	KnownNodes      = []string{"localhost:3000"}
	blocksInTransit = [][]byte{}
	memoryPool      = make(map[string]blockchain.Transaction)
)

// Addr contains the address of a node
type Addr struct {
	AddrList []string
}

// Block contains a block
type Block struct {
	AddrFrom string
	Block    []byte
}

// GetBlocks contains the address of a node
type GetBlocks struct {
	AddrFrom string
}

// GetData contains the address of a node and the hash of a block
type GetData struct {
	AddrFrom string
	Type     string
	ID       []byte
}

// Inventory contains the address of a node and the hash of a block
type Inventory struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

// Tx contains a transaction
type Tx struct {
	AddrFrom    string
	Transaction []byte
}

// Version contains the best height of a node and the address of a node
type Version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

// CmdToBytes converts a string command to a byte array
func CmdToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

// BytesToCmd converts a byte array to a string command
func BytesToCmd(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

func CloseDB(chain *blockchain.BlockChain) {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	d.WaitForDeathWithFunc(func() {
		chain.Database.Close()
	})
}
