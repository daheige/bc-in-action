package infras

// Blockchain 区块链数据
type Blockchain struct {
	Blocks []*Block // 区块链数据，这里使用切片进行模拟，实际上区块链的数据是一个分布式账本，分布式数据库的方式存储数据
}

// NewBlockchain 创建一个blockchain区块链
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

// AddBlock 添加区块
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}
