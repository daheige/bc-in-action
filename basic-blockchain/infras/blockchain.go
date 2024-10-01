package infras

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// Blockchain 区块链数据
type Blockchain struct {
	tip []byte // 存储数据

	// 区块链数据，这里使用切片进行模拟，实际上区块链的数据是一个分布式账本，分布式数据库的方式存储数据
	db           *bolt.DB // 数据库使用bolt db
	blocksBucket string
}

const (
	defaultDBFile       = "blockchain.db"
	defaultBlocksBucket = "blocks"
)

// NewBlockchain 创建一个blockchain区块链
func NewBlockchain(dbFile, blocksBucket string) (*Blockchain, error) {
	if dbFile == "" {
		dbFile = defaultDBFile
	}
	if blocksBucket == "" {
		blocksBucket = defaultBlocksBucket
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			// 如果不存在，就生成创世块，创建 bucket，并将区块保存到里面，然后更新 l 键以存储链中最后一个块的哈希
			fmt.Println("No existing blockchain found. Creating a new blockchain...")
			genesis := NewGenesisBlock() // 创建创世区块
			b, err = tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				return err
			}

			var data []byte
			data, err = genesis.Serialize()
			if err != nil {
				return err
			}

			err = b.Put([]byte(genesis.Hash), data)
			if err != nil {
				return err
			}

			err = b.Put([]byte("l"), []byte(genesis.Hash))
			if err != nil {
				return err
			}
			tip = []byte(genesis.Hash)
		} else { // 如果存在，就从中读取 l 键
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := &Blockchain{tip, db, blocksBucket}
	return bc, nil
}

// AddBlock 添加区块
func (bc *Blockchain) AddBlock(data string) error {
	var lastHash []byte
	// 这是 BoltDB 事务的另一个类型（只读）。
	// 在这里，我们会从数据库中获取最后一个块的哈希，然后用它来挖出一个新的块的哈希
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bc.blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		return err
	}

	// 新的block
	newBlock := NewBlock(data, string(lastHash))
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bc.blocksBucket))
		d, execErr := newBlock.Serialize()
		if execErr != nil {
			return execErr
		}

		h := []byte(newBlock.Hash)
		execErr = b.Put(h, d)
		if execErr != nil {
			return execErr
		}

		execErr = b.Put([]byte("l"), h)
		bc.tip = h
		return nil
	})

	return err
}

// BlockchainIterator 区块链迭代器
// 每当要对链中的块进行迭代时，我们就会创建一个迭代器，
// 里面存储了当前迭代的块哈希（currentHash）和数据库的连接（db）
type BlockchainIterator struct {
	currentHash  []byte
	db           *bolt.DB
	blocksBucket string
}

// Iterator 迭代器
// 注意，迭代器的初始状态为链中的 tip，因此区块将从尾到头（创世块为头），也就是从最新的到最旧的进行获取。
// 实际上，选择一个 tip 就是意味着给一条链“投票”。一条链可能有多个分支，最长的那条链会被认为是主分支。
// 在获得一个 tip （可以是链中的任意一个块）之后，我们就可以重新构造整条链，找到它的长度和需要构建它的工作。
// 这同样也意味着，一个 tip 也就是区块链的一种标识符。
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{currentHash: bc.tip, db: bc.db, blocksBucket: bc.blocksBucket}
	return bci
}

// Next 迭代block
func (i *BlockchainIterator) Next() (*Block, error) {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(i.blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		var err error
		block, err = DeserializeBlock(encodedBlock)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	i.currentHash = []byte(block.PrevHash)
	return block, err
}
