package blockchain

import (
	"log"

	"github.com/dgraph-io/badger"
)

type blockChainIterator struct {
	currentHash []byte
	Database    *badger.DB
}

func (chain *BlockChain) Iterator() *blockChainIterator {
	iter := &blockChainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *blockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.currentHash)
		if err != nil {
			log.Panic(err)
		}

		var encodedBlock []byte
		err = item.Value(func(val []byte) error {
			encodedBlock = val
			return nil
		})
		block = Deserialize(encodedBlock)

		return err
	})

	if err != nil {
		log.Panic(err)
	}

	iter.currentHash = block.PrevHash

	return block
}
