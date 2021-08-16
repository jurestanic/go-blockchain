package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

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

func genesis() *block {
	return createBlock("Genesis Block", []byte{})
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
