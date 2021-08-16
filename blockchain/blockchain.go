package blockchain

import (
	"log"

	"github.com/dgraph-io/badger"
)

const dbPath = "./tmp/blocks"

type BlockChain struct {
	lastHash []byte
	Database *badger.DB
}

type blockChainIterator struct {
	currentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			gen := genesis()
			err = txn.Set(gen.Hash, gen.serialize())
			if err != nil {
				log.Panic(err)
			}
			err = txn.Set([]byte("lh"), gen.Hash)
			lastHash = gen.Hash

			return err
		}

		item, err := txn.Get([]byte("lh"))
		if err != nil {
			log.Panic(err)
		}

		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})

		return err
	})

	if err != nil {
		log.Panic(err)
	}

	blockchain := BlockChain{lastHash: lastHash, Database: db}

	return &blockchain
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			log.Panic(err)
		}
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})

		return err
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := createBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.serialize())
		if err != nil {
			log.Panic(err)
		}
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.lastHash = newBlock.Hash

		return err
	})

	if err != nil {
		log.Panic(err)
	}
}

func (chain *BlockChain) Iterator() *blockChainIterator {
	iter := &blockChainIterator{chain.lastHash, chain.Database}

	return iter
}

func (iter *blockChainIterator) Next() *block {
	var block *block

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
		block = deserialize(encodedBlock)

		return err
	})

	if err != nil {
		log.Panic(err)
	}

	iter.currentHash = block.PrevHash

	return block
}
