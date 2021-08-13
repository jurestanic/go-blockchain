package blockchain

import (
	"fmt"
)

type blockChain struct {
	blocks []*block
}

type block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func createBlock(data string, prevHash []byte) *block {
	block := &block{Hash: []byte{}, Data: []byte(data), PrevHash: prevHash, Nonce: 0}
	pow := newProof(block)
	nonce, hash := pow.run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

func (chain *blockChain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	newBlock := createBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, newBlock)
}

func genesis() *block {
	return createBlock("Genesis Block", []byte{})
}

func InitBlockChain() *blockChain {
	return &blockChain{[]*block{genesis()}}
}

func PrintBlockChain(chain *blockChain) {
	for _, block := range chain.blocks {
		fmt.Printf("Prev Hash: %x\n", block.PrevHash)
		fmt.Printf("Data In Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
