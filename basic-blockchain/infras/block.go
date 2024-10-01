package infras

import (
	"encoding/hex"
	"time"
)

// Block 区块信息，基本字段
type Block struct {
	Timestamp int64  // 当前时间戳
	Data      string // 当前区块存储的实际有效信息，也就是交易
	PrevHash  string // 前一个块的哈希值

	// hash格式 0000002aa0fcc49e60d8ba1d9061c863c656e51f7fad190566b322034992e4ce
	Hash  string // 当前块的哈希值
	Nonce int64  // 计数器数字
}

// NewGenesisBlock 创建创世区块 genesis block，也就是区块链中的第一个block块
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", "0")
}

// NewBlock 创建一个区块block
func NewBlock(data string, prevHash string) *Block {
	block := &Block{
		Timestamp: time.Now().Unix(),
		Data:      data,
		PrevHash:  prevHash,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run() // 计算hash值

	block.Hash = hex.EncodeToString(hash)
	block.Nonce = nonce

	return block
}
