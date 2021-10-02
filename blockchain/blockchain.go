package blockchain

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "Coinbase Transaction From Genesis"
)

type BlockChain struct {
	lastHash []byte
	Database *badger.DB
}

type blockChainIterator struct {
	currentHash []byte
	Database    *badger.DB
}

func InitBlockChain(address string) *BlockChain {
	var lastHash []byte

	if DBexists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(address, genesisData)
		gen := genesis(cbtx)

		err = txn.Set(gen.Hash, gen.serialize())
		if err != nil {
			log.Panic(err)
		}
		err = txn.Set([]byte("lh"), gen.Hash)
		lastHash = gen.Hash

		return err
	})

	if err != nil {
		log.Panic(err)
	}

	blockchain := BlockChain{lastHash: lastHash, Database: db}

	return &blockchain
}

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func ContinueBlockChain(address string) *BlockChain {
	if DBexists() == false {
		fmt.Println("No existing blockchain found!")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
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

func (chain *BlockChain) AddBlock(transactions []*Transaction) {
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

	newBlock := createBlock(transactions, lastHash)

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

func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var utxo []Transaction

	spentTXO := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transaction {
			TxID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXO[TxID] != nil {
					for _, spentOut := range spentTXO[TxID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					utxo = append(utxo, *tx)
				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXO[inTxID] = append(spentTXO[inTxID], in.Out)
					}
				}
			}

		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return utxo
}

func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXO []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXO = append(UTXO, out)
			}
		}
	}

	return UTXO
}

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}
