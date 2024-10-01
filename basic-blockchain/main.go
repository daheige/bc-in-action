package main

import (
	"fmt"

	"basic-blockchain/infras"
)

func main() {
	// 创建一个区块链对象
	bc, err := infras.NewBlockchain("blockchain.db", "blocks")
	if err != nil {
		fmt.Println("new blockchain err: ", err)
		return
	}
	defer bc.Close()

	cli := infras.NewCLI(bc)
	cli.Run()
}
