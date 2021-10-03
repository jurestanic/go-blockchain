package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type block struct {
	Hash        []byte
	Transaction []*Transaction
	PrevHash    []byte
	Nonce       int
}

func (b *block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transaction {
		txHashes = append(txHashes, tx.Serialize())
	}

	tree := NewMerkleTree(txHashes)

	return tree.RootNode.Data
}

func createBlock(txs []*Transaction, prevHash []byte) *block {
	block := &block{Hash: []byte{}, Transaction: txs, PrevHash: prevHash, Nonce: 0}
	pow := newProof(block)
	nonce, hash := pow.run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

func genesis(coinbase *Transaction) *block {
	return createBlock([]*Transaction{coinbase}, []byte{})
}

func (b *block) serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	if err != nil {
		log.Panic(err)
	}

	return res.Bytes()
}

func deserialize(data []byte) *block {
	var block block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}

	return &block
}
