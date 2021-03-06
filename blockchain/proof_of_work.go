package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const difficulty = 12

type proofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProof(b *Block) *proofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	pow := &proofOfWork{b, target}

	return pow
}

func (pow *proofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.initData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.target) == -1 {
			// less than target (we did the work)!
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

func (pow *proofOfWork) initData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.PrevHash,
		pow.block.HashTransactions(),
		intToBytes(int64(nonce)),
		intToBytes(int64(difficulty)),
	}, []byte{})

	return data
}

func intToBytes(num int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, num)
	if err != nil {
		// TODO: Create app-specific logger!
		log.Panic(err)
	}

	return buffer.Bytes()
}
