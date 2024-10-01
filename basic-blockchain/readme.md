# 使用go实现一个简单的区块链封装
```go
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 实现一个简单的区块链

// Block 区块信息，基本字段
type Block struct {
	Timestamp int64  // 当前时间戳
	Data      string // 当前区块存储的实际有效信息，也就是交易
	PrevHash  string // 前一个块的哈希值
	Hash      string // 当前块的哈希值
}

// NewBlock 创建block
func NewBlock(data string, prevHash string) *Block {
	block := &Block{
		Timestamp: time.Now().Unix(),
		Data:      data,
		PrevHash:  prevHash,
	}

	block.setHash()
	return block
}

// Blockchain 区块链数据
type Blockchain struct {
	blocks []*Block // 区块链数据，这里使用切片进行模拟，实际上区块链的数据是一个分布式账本，分布式数据库的方式存储数据
}

// AddBlock 添加区块
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

// 创建创世区块 genesis block，也就是区块链中的第一个block块
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", "0")
}

// NewBlockchain 创建一个blockchain区块链
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

// setHash 设置hash值
func (b *Block) setHash() {
	timestamp := strconv.FormatInt(b.Timestamp, 10)
	headers := strings.Join([]string{
		b.PrevHash, b.Data, timestamp,
	}, "")
	hash := sha256.Sum256([]byte(headers))
	b.Hash = hex.EncodeToString(hash[:])
}

func main() {
	// 测试区块链是否正常运行
	// 创建一个区块链对象
	bc := NewBlockchain()

	// 添加区块
	bc.AddBlock("send 1 btc to van")
	bc.AddBlock("send 2 btc to alex")
	bc.AddBlock("send 3 more btc to van")
	for _, block := range bc.blocks {
		fmt.Println("current block info")
		fmt.Printf("Previous hash: %s\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Println("")
	}
}
```
输出结果如下：
```
current block info
Previous hash: 0
Data: Genesis Block
Hash: 99eb04b2c71a7cb620a3d9a384dd5742e37753a82c15c07990a20d115fdaaacb

current block info
Previous hash: 99eb04b2c71a7cb620a3d9a384dd5742e37753a82c15c07990a20d115fdaaacb
Data: send 1 btc to van
Hash: 379853060647d657105f42706cd1e47a4fc12c5a268ff424ca4cff6e3b8fae68

current block info
Previous hash: 379853060647d657105f42706cd1e47a4fc12c5a268ff424ca4cff6e3b8fae68
Data: send 2 btc to alex
Hash: 968151cd871ba62a5ed4dbf0192c80e0bad21b4d3b0f4129f1f00a45432696a3

current block info
Previous hash: 968151cd871ba62a5ed4dbf0192c80e0bad21b4d3b0f4129f1f00a45432696a3
Data: send 3 more btc to van
Hash: 9060bf912a291229994d148eebd9cff378467964ba041a717fd788e68f34df4d
```

# 工作量证明POW
- 区块链的一个关键点就是，一个人必须经过一系列困难的工作，才能将数据放入到区块链中。正是由于这种困难的工作，才保证了区块链的安全和一致。此外，完成这个工作的人，也会获得相应奖励（这也就是通过挖矿获得币）。
- 这个机制与生活现象非常类似：一个人必须通过努力工作，才能够获得回报或者奖励，用以支撑他们的生活。在区块链中，是通过网络中的参与者（矿工）不断的工作来支撑起了整个网络。矿工不断地向区块链中加入新块，然后获得相应的奖励。在这种机制的作用下，新生成的区块能够被安全地加入到区块链中，它维护了整个区块链数据库的稳定性。值得注意的是，完成了这个工作的人必须要证明这一点，即他必须要证明他的确完成了这些工作。
- 整个 “努力工作并进行证明” 的机制，就叫做工作量证明（proof-of-work）。要想完成工作非常地不容易，因为这需要大量的计算能力：即便是高性能计算机，也无法在短时间内快速完成。另外，这个工作的困难度会随着时间不断增长，以保持每 10 分钟出 1 个新块的速度。在比特币中，这个工作就是找到一个块的哈希，同时这个哈希满足了一些必要条件。这个哈希，也就充当了证明的角色。因此，寻求证明（寻找有效哈希），就是矿工实际要做的事情。

# 哈希计算
获得指定数据的一个哈希值的过程，就叫做哈希计算。 一个哈希，就是对所计算数据的一个唯一表示。对于一个哈希函数，输入任意大小的数据，它会输出一个固定大小的哈希值。下面是哈希的几个关键特性：
1. 无法从一个哈希值恢复原始数据。也就是说，哈希并不是加密。
2. 对于特定的数据，只能有一个哈希，并且这个哈希是唯一的。
3. 即使是仅仅改变输入数据中的一个字节，也会导致输出一个完全不同的哈希。

- 哈希函数被广泛用于检测数据的一致性。软件提供者常常在除了提供软件包以外，还会发布校验和。当下载完一个文件以后，你可以用哈希函数对下载好的文件计算一个哈希，并与作者提供的哈希进行比较，以此来保证文件下载的完整性。
- 在区块链中，哈希被用于保证一个块的一致性。哈希算法的输入数据包含了前一个块的哈希，因此使得不太可能（或者，至少很困难）去修改链中的一个块：因为如果一个人想要修改前面一个块的哈希，那么他必须要重新计算这个块以及后面所有块的哈希。

# Hashcash
比特币使用Hashcash，它是一个最初用来防止垃圾邮件的工作量证明算法，可以分解为如下步骤：
1. 获取一些公开的数据（比如，如果是一个email的话，它可以是接收者的邮件地址；在比特币中，它是区块头）。
2. 给这个公开数据添加一个计数器counter。计数器默认从0开始。
3. 将data(数据)和counter(计数器)组合到一起，获得一个哈希。
4. 检查这个哈希是否符合一定的条件：
   - 如果符合条件，结束。
   - 如果不符合条件，增加计数器，重复3-4

因此，这是一个暴力算法：改变计数器，计算新的哈希，检查，增加计数器，计算新的哈希，检查...，如此往复。计算成本是非常高的。

- 一个哈希要满足的必要条件。在原始的 Hashcash 实现中，它的要求是 “一个哈希的前 20 位必须是 0”。
- 在比特币中，这个要求会随着时间而不断变化。因为按照设计，必须保证每 10 分钟生成一个块，而不论计算能力会随着时间增长，或者是会有越来越多的矿工进入网络，所以需要动态调整这个必要条件。

# 工作量实现
```go
// 工作量证明
 // 你可以把目标想象为一个范围的上界：如果一个数（由哈希转换而来）比上界要小，那么是有效的，反之无效。
 // 因为要求比上界要小，所以会导致有效数字并不会很多。
 // 因此，也就需要通过一些困难的工作（一系列反复地计算），才能找到一个有效的数字。
 data1 := []byte("i like music")
 data2 := []byte("i like programing")
 target := big.NewInt(1)
 target.Lsh(target, uint(256-infras.TargetBits))
 fmt.Printf("%x\n", sha256.Sum256(data1))
 fmt.Printf("%64x\n", target)
 fmt.Printf("%x\n", sha256.Sum256(data2))
```
具体实现，参考 proof_of_work.go ，运行main.go结果如下：
```
Mining the block containing "Genesis Block"
000000be2d394baeb11921eec362947b0fc24dcedea3c07026a79b66e1808104
Mining the block containing "send 1 btc to van"
0000008429a24d70a84ba30aa26f58c8fb58a355513228b2b215b4c42ce3c757
Mining the block containing "send 2 btc to alex"
0000002aa0fcc49e60d8ba1d9061c863c656e51f7fad190566b322034992e4ce
Mining the block containing "send 3 more btc to van"
00000050415ef34244531acb751c191c1fc9f2197cd6b02e431fbc59738fc8dc
current block info
Previous hash: 0
Data: Genesis Block
Hash: 000000be2d394baeb11921eec362947b0fc24dcedea3c07026a79b66e1808104

current block info
Previous hash: 000000be2d394baeb11921eec362947b0fc24dcedea3c07026a79b66e1808104
Data: send 1 btc to van
Hash: 0000008429a24d70a84ba30aa26f58c8fb58a355513228b2b215b4c42ce3c757

current block info
Previous hash: 0000008429a24d70a84ba30aa26f58c8fb58a355513228b2b215b4c42ce3c757
Data: send 2 btc to alex
Hash: 0000002aa0fcc49e60d8ba1d9061c863c656e51f7fad190566b322034992e4ce

current block info
Previous hash: 0000002aa0fcc49e60d8ba1d9061c863c656e51f7fad190566b322034992e4ce
Data: send 3 more btc to van
Hash: 00000050415ef34244531acb751c191c1fc9f2197cd6b02e431fbc59738fc8dc
```
成功了！你可以看到每个哈希都是 3 个字节的 0 开始，并且获得这些哈希需要花费一些时间。

还剩下一件事情需要做，对工作量证明进行验证：
```go
// 对工作量进行验证
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	// pow.block.Nonce 是 Run 方法的执行结果，计数器
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
```
运行结果如下：
```
Mining the block containing "Genesis Block"
hash: 00000089bb31612446ebc893c8de5a28a0f005fe5e5e136a65b3fa2906fd0aab

Mining the block containing "send 1 btc to van"
hash: 00000038732531a7dc823f34ac51de04508a7a03c6ef38ab34238baf025aec31

Mining the block containing "send 2 btc to alex"
hash: 0000000dd0f3175c17b6bdef46328c0ea4d73de65a2e5983a1add5bd37176625

Mining the block containing "send 3 more btc to van"
hash: 000000c9aa6c7fd413939bc299daebcd34925768182ed73e411e601f1928c97e

Previous hash: 0
Data: Genesis Block
Hash: 00000089bb31612446ebc893c8de5a28a0f005fe5e5e136a65b3fa2906fd0aab
pow validate begin...
PoW: true
pow validate end

Previous hash: 00000089bb31612446ebc893c8de5a28a0f005fe5e5e136a65b3fa2906fd0aab
Data: send 1 btc to van
Hash: 00000038732531a7dc823f34ac51de04508a7a03c6ef38ab34238baf025aec31
pow validate begin...
PoW: true
pow validate end

Previous hash: 00000038732531a7dc823f34ac51de04508a7a03c6ef38ab34238baf025aec31
Data: send 2 btc to alex
Hash: 0000000dd0f3175c17b6bdef46328c0ea4d73de65a2e5983a1add5bd37176625
pow validate begin...
PoW: true
pow validate end

Previous hash: 0000000dd0f3175c17b6bdef46328c0ea4d73de65a2e5983a1add5bd37176625
Data: send 3 more btc to van
Hash: 000000c9aa6c7fd413939bc299daebcd34925768182ed73e411e601f1928c97e
pow validate begin...
PoW: true
pow validate end
```
从输出结果看出，产生3个块，花了2分钟以上，比没有工作量证明慢很多了，也就是成本高出了很多。

- 我们离真正的区块链又进了一步：现在需要经过一些困难的工作才能加入新的块，因此挖矿就有可能了。
- 但是，它仍然缺少一些至关重要的特性：区块链数据库并不是持久化的，没有钱包，地址，交易，也没有共识机制。

到目前为止，我们已经构建了一个有工作量证明机制的区块链。有了工作量证明，挖矿也就有了着落。 虽然目前距离一个有着完整功能的区块链越来越近了，但是它仍然缺少了一些重要的特性。

我们会将区块链持久化到一个数据库中，然后会提供一个简单的命令行接口，用来完成一些与区块链的交互操作。
本质上，区块链是一个分布式数据库，不过，我们暂时先忽略 “分布式” 这个部分，后面会详细说明这些内容。

在这个例子中，在每次运行程序时，简单地将区块链存储在内存中。那么一旦程序退出，所有的内容就都消失了。我们没有办法再次使用这条链，也没有办法与其他人共享，所以我们需要把它存储到磁盘或者分布式数据库上。

那么，我们要用哪个数据库呢？实际上，任何一个数据库都可以。在 比特币原始论文 中，并没有提到要使用哪一个具体的数据库，它完全取决于开发者如何选择。 Bitcoin Core ，最初由中本聪发布，现在是比特币的一个参考实现，它使用的是 LevelDB。

# BoltDB
具有以下特征：
1. 简单
2. 使用go实现
3. 不需要运行一个服务器
4. 能够允许开发者构造想要的数据结构

blot官方介绍：https://github.com/boltdb/bolt
- Bolt 是一个纯键值存储的 Go 数据库，启发自 Howard Chu 的 LMDB. 它旨在为那些无须一个像 Postgres 和 MySQL 这样有着完整数据库服务器的项目，提供一个简单，快速和可靠的数据库。
- 由于 Bolt 意在用于提供一些底层功能，简洁便成为其关键所在。它的 API 并不多，并且仅关注值的获取和设置。仅此而已。

听起来跟我们的需求完美契合！来快速过一下：

- Bolt 使用键值存储，这意味着它没有像 SQL RDBMS （MySQL，PostgreSQL 等等）的表，没有行和列。相反，数据被存储为键值对（key-value pair，就像 Golang 的 map）。键值对被存储在 bucket 中，这是为了将相似的键值对进行分组（类似 RDBMS 中的表格）。
- 因此，为了获取一个值，你需要知道一个 bucket 和一个键（key）。
- 需要注意的一个事情是，Bolt 数据库没有数据类型：键和值都是字节数组（byte array）。鉴于需要在里面存储 Go 的结构（准确来说，也就是存储Block（块）），我们需要对它们进行序列化，也就说，实现一个从 Go struct 转换到一个 byte array 的机制，同时还可以从一个 byte array 再转换回 Go struct。
- 虽然我们将会使用 encoding/gob 来完成这一目标，但实际上也可以选择使用 JSON, XML, Protocol Buffers 等等。之所以选择使用 encoding/gob, 是因为它很简单，而且是 Go 标准库的一部分。

虽然blot不再活跃，但是有活跃的fork版本：https://github.com/etcd-io/bbolt

# 数据库结构
在开始实现持久化的逻辑之前，我们首先需要决定到底要如何在数据库中进行存储。为此，我们可以参考 Bitcoin Core 的做法：

简单来说，Bitcoin Core 使用两个 “bucket” 来存储数据：
- 其中一个 bucket 是 blocks，它存储了描述一条链中所有块的元数据
- 另一个 bucket 是 chainstate，存储了一条链的状态，也就是当前所有的未花费的交易输出，和一些元数据

此外，出于性能的考虑，Bitcoin Core 将每个区块（block）存储为磁盘上的不同文件。如此一来，就不需要仅仅为了读取一个单一的块而将所有（或者部分）的块都加载到内存中。

# 序列化处理
```go
// Serialize 序列化数据
func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

// DeserializeBlock 反序列化数据
func DeserializeBlock(data []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}
```

# 持久化
使用 bolt "go.etcd.io/bbolt" 实现

# 命令行接口
```shell
make build && make test
```
输出结果：
```
====> Go build
./basic-blockchain printchain
No existing blockchain found. Creating a new blockchain...
Mining the block containing "Genesis Block"
hash: 00000044a7514c43349332f30edddc610429f3b2031e7f7bda8fd1f82b7147c4

Previous hash: 
Data: Genesis Block
Hash: 00000044a7514c43349332f30edddc610429f3b2031e7f7bda8fd1f82b7147c4
pow validate begin...
PoW: true
pow validate end

./basic-blockchain addblock -data "send 1 btc to van"
Mining the block containing "send 1 btc to van"
hash: 000000170c6c948b4469ad0b7455e7070d63ed94e0a716a30081db84c3d4b381

./basic-blockchain addblock -data "send 2 btc to alex"
Mining the block containing "send 2 btc to alex"
hash: 000000734cda8779ba25d2e921d7a0dcc7acaaaed8701dd1d40d1d9383358d83

./basic-blockchain addblock -data "send 3 more btc to van"
Mining the block containing "send 3 more btc to van"
hash: 000000e6072d67c2cfc1e7703b6fc96605275eaedac355ad11bede35fb1f4c61

./basic-blockchain printchain
Previous hash: 000000734cda8779ba25d2e921d7a0dcc7acaaaed8701dd1d40d1d9383358d83
Data: send 3 more btc to van
Hash: 000000e6072d67c2cfc1e7703b6fc96605275eaedac355ad11bede35fb1f4c61
pow validate begin...
PoW: true
pow validate end

Previous hash: 000000170c6c948b4469ad0b7455e7070d63ed94e0a716a30081db84c3d4b381
Data: send 2 btc to alex
Hash: 000000734cda8779ba25d2e921d7a0dcc7acaaaed8701dd1d40d1d9383358d83
pow validate begin...
PoW: true
pow validate end

Previous hash: 00000044a7514c43349332f30edddc610429f3b2031e7f7bda8fd1f82b7147c4
Data: send 1 btc to van
Hash: 000000170c6c948b4469ad0b7455e7070d63ed94e0a716a30081db84c3d4b381
pow validate begin...
PoW: true
pow validate end

Previous hash: 
Data: Genesis Block
Hash: 00000044a7514c43349332f30edddc610429f3b2031e7f7bda8fd1f82b7147c4
pow validate begin...
PoW: true
pow validate end
```
