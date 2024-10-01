package infras

import (
	"flag"
	"fmt"
	"os"
)

// CLI 命令cli
type CLI struct {
	bc *Blockchain
}

// NewCLI 创建一个区块链cli工具对象
func NewCLI(bc *Blockchain) *CLI {
	return &CLI{bc: bc}
}

// Run 运行命令行
func (cli *CLI) Run() {
	cli.validateArgs() // 验证参数

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printBlockCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// 给 addblock 添加 -data 标志
	addBlockData := addBlockCmd.String("data", "", "block data")
	var err error
	switch os.Args[1] { // 判断第2个参数
	case "addblock":
		err = addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		err = printBlockCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(0)
	}

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(0)
	}

	// 解析相关的 flag 子命令
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(0)
		}

		err = cli.addBlock(*addBlockData)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}

	if printBlockCmd.Parsed() {
		cli.printChain()
	}
}

const usage = `
	Usage: 
		addblock -data BLOCK_DATA add a block to the blockchain
		printchain 				  print all the blocks of the blockchain
`

func (cli *CLI) printUsage() {
	fmt.Println(usage)
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(0)
	}
}

func (cli *CLI) addBlock(data string) error {
	return cli.bc.AddBlock(data)
}

func (cli *CLI) printChain() {
	// BlockchainIterator 对区块链中的区块进行迭代
	bci := cli.bc.Iterator()
	for {
		block, err := bci.Next()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Printf("Previous hash: %s\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %s\n", block.Hash)

		// 验证工作量是否有效
		fmt.Println("pow validate begin...")
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %v\n", pow.Validate())
		fmt.Println("pow validate end")
		fmt.Println("")

		if len(block.PrevHash) == 0 {
			break
		}
	}
}
