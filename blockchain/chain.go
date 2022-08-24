package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/nbright/nomadcoin/db"
	"github.com/nbright/nomadcoin/utils"
)

// func (b *BlockChain) listBlocks() {
// 	for _, block := range b.blocks {
// 		fmt.Printf("Data: %s\n", block.data)
// 		fmt.Printf("Hash: %s\n", block.hash)
// 		fmt.Printf("PreHash: %s\n", block.prevHash)
// 	}
// }

type blockChain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockChain
var once sync.Once

func (b *blockChain) restore(data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	decoder.Decode(b)
}

func (b *blockChain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockChain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func BlockChain() *blockChain {
	if b == nil {
		once.Do(func() {
			b = &blockChain{"", 0}
			fmt.Printf("1NewestHash: %s\nHeight:%d", b.NewestHash, b.Height)
			// search for "checkpoint" on the db
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				fmt.Printf("Restoring...")
				// restore b from bytes
				b.restore(checkpoint)
			}

		})
	}
	fmt.Printf("3NewestHash: %s\nHeight:%d", b.NewestHash, b.Height)
	return b
}
