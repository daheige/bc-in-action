package infras

import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// ProofOfWork pow工作量证明
// 里面存储了指向一个块(block)和一个目标(target)的指针，计数器
// 这里使用大整数
// 在检查的时候，它会将哈希与目标进行比较：先把哈希转换为一个大整数，然后检查它是否小于目标。
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// TargetBits 挖矿的难度值
// 在比特币中，当一个块被挖出来以后，“target bits” 代表了区块头里存储的难度，也就是开头有多少个 0。
// 这里的 24 指的是算出来的哈希前 24 位必须是 0，如果用 16 进制表示，就是前 6 位必须是 0，
// 这一点从最后的输出可以看出来。目前我们并不会实现一个动态调整目标的算法，所以将难度定义为一个全局的常量即可。
// 24 其实是一个可以任意取的数字，其目的只是为了有一个目标（target）而已，这个目标占据不到 256 位的内存空间。
// 同时，我们想要有足够的差异性，但是又不至于大的过分，因为差异性越大，就越难找到一个合适的哈希。
const TargetBits = 24

// NewProofOfWork 创建一个工作证明对象
// 在 NewProofOfWork 函数中，我们将 big.Int 初始化为 1，然后左移 256 - targetBits 位。
// 256 是一个 SHA-256 哈希的位数，
// 我们将要使用的是 SHA-256 哈希算法。target（目标） 的 16 进制形式为：
// 0x10000000000000000000000000000000000000000000000000000000000
// 它在内存上占据了 29 个字节
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1) // 初始化为1
	// >> 右移动 256-targetBits 位
	// 256是一个SHA-256哈希的位数，因为block是用的SHA-256算法，target可以表示为一个16进制的数字
	target.Lsh(target, uint(256-TargetBits))
	pow := &ProofOfWork{block, target}

	return pow
}

// 准备数据
// 这里的 nonce，就是上面 Hashcash 所提到的计数器，它是一个密码学术语。
// 返回的数据是:prevHash+data+timestamp+targetBits+nonce
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := strings.Join([]string{
		pow.block.PrevHash,
		pow.block.Data,
		strconv.FormatInt(pow.block.Timestamp, 10),
		strconv.FormatInt(TargetBits, 10),
		strconv.FormatInt(nonce, 10),
	}, "")

	return []byte(data)
}

// Run 实现 PoW 算法的核心
func (pow *ProofOfWork) Run() (int64, []byte) {
	var (
		hashInt big.Int // hash的整形表示
		hash    [32]byte
		nonce   int64 // 计数器
	)

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce) // 准备数据

		// 用 SHA-256 对数据进行哈希
		//  hash = SHA256(prevHash+data+timestamp+targetBits+nonce)
		hash = sha256.Sum256(data)

		hashInt.SetBytes(hash[:]) // 将哈希换成一个大整数

		// 如果哈希值小于pow.target
		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("hash: %x\n", hash)
			break
		} else {
			nonce += 1
		}
	}

	fmt.Println("")

	return nonce, hash[:]
}

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
