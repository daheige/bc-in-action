package main

import (
	"fmt"

	"basic-blockchain/infras"
)

func main() {
	// 测试区块链是否正常运行
	// 创建一个区块链对象
	bc := infras.NewBlockchain()

	// 添加区块
	bc.AddBlock("send 1 btc to van")
	bc.AddBlock("send 2 btc to alex")
	bc.AddBlock("send 3 more btc to van")
	for _, block := range bc.Blocks {
		fmt.Printf("Previous hash: %s\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %s\n", block.Hash)

		// 验证工作量是否有效
		fmt.Println("pow validate begin...")
		pow := infras.NewProofOfWork(block)
		fmt.Printf("PoW: %v\n", pow.Validate())
		fmt.Println("pow validate end")
		fmt.Println("")
	}
}
