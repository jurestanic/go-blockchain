package main

import (
	"github.com/jurestanic/go-blockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	chain.AddBlock("First Block !")
	chain.AddBlock("Second Block !")
	chain.AddBlock("Third Block !")

	blockchain.PrintBlockChain(chain)
}
