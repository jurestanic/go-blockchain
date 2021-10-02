package main

import (
	"os"

	"github.com/jurestanic/go-blockchain/cli"
)

func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()
}
