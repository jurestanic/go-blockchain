package main

import (
	"os"

	"github.com/jurestanic/go-blockchain/blockchain"
	"github.com/jurestanic/go-blockchain/cli"
)

func main() {
	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	cli := cli.CommandLine{Blockchain: chain}
	cli.Run()
}
