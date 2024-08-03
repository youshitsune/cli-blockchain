package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

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
	getDataCmd := flag.NewFlagSet("getdata", flag.ExitOnError)

	addBlockPath := addBlockCmd.String("f", "", "File with data")
	addBlockName := addBlockCmd.String("n", "", "Name of block")
	getDataData := getDataCmd.String("hash", "", "Hash of block")

	switch os.Args[1] {
	case "addblock":
		addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		printChainCmd.Parse(os.Args[2:])
	case "getdata":
		getDataCmd.Parse(os.Args[2:])
	default:
		cli.PrintUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockPath == "" && *addBlockName == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockPath, *addBlockName)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if getDataCmd.Parsed() {
		if *getDataData == "" {
			getDataCmd.Usage()
			os.Exit(1)
		}
		cli.getData(*getDataData)
	}

}

func (cli *CLI) addBlock(path, name string) {
	res := cli.Req("addblock", "req")
	ctx, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	prevHash, _ := hex.DecodeString(res)
	block := NewBlock(name, string(ctx), prevHash)
	res = cli.Req("addblock", string(block.Serialize()))
	fmt.Println(res)
}

func (cli *CLI) printChain() {
	res := cli.Req("", "")
	blocks := strings.Split(res, "|/")
	fmt.Println(len(blocks))
	fmt.Println()
	for i := range blocks {
		block := Deserialize([]byte(blocks[i]))
		fmt.Printf("Name: %s\n", block.Name)
		fmt.Printf("Data: %d\n", len(block.Data))
		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := NewPoW(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}

func (cli *CLI) getData(hash string) {
	res := strings.Split(cli.Req("getdata", hash), "|/")
	data := []byte(res[0])
	name := res[1]

	err := os.WriteFile(name, data, 0677)
	if err != nil {
		fmt.Println(err)
	}
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
