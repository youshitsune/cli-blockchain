package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/go-zoox/fetch"
)

type CLI struct {
	url string
}

func (cli *CLI) Req(url, data string) string {
	var response *fetch.Response
	var err error
	if len(data) > 0 {
		req := fetch.Query{}
		req.Set("data", data)
		response, err = fetch.Post(cli.url+url, &fetch.Config{Query: req})
		if err != nil {
			fmt.Printf("Can not reach the server: %v", err)
			os.Exit(1)
		}

	} else {
		response, err = fetch.Get(cli.url + url)
		if err != nil {
			fmt.Printf("Can not reach the server: %v", err)
			os.Exit(1)
		}
	}

	return response.String()

}

func (cli *CLI) Run() {

	cli.ValidateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addblock":
		addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		printChainCmd.Parse(os.Args[2:])
	default:
		cli.PrintUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}

func (cli *CLI) addBlock(data string) {
	res := cli.Req("addblock", "req")
	prevHash, _ := hex.DecodeString(res)
	block := NewBlock(data, prevHash)
	res = cli.Req("addblock", block.Serialize())
	fmt.Println(res)
}

func (cli *CLI) printChain() {
	res := cli.Req("", "")
	fmt.Println(res)
}

func (cli *CLI) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}
}

func (cli *CLI) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("\taddblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("\tprintchain - print all the blocks of the blockchain")
}

func main() {
	cli := CLI{"http://localhost:8080/"}
	cli.Run()
}
