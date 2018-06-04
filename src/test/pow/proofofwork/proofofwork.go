package proofofwork

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"blockChainMorp/utils"
	"blockChainMorp/types"
	"blockChainMorp/transaction/transactions"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 4 * 5

type proofOfWork struct {
	block  *types.Block
	target *big.Int
}

// NewProofOfWork builds and returns a new ProofOfWork
func NewProofOfWork() *proofOfWork {
	pow := new(proofOfWork)

	return pow
}

// InitBlock inits a block in the proof-of-work
func (pow *proofOfWork) InitBlock(b *types.Block) {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow.target = target
	pow.block = b
}

func (pow *proofOfWork) prepareData(nonce int) []byte {
	itx := transactions.NewTrans()
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			itx.HashTransactions(pow.block),
			utils.IntToHex(pow.block.Timestamp),
			utils.IntToHex(int64(targetBits)),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// Run performs a proof-of-work
func (pow *proofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Println("Mining a new block...")
	for nonce < maxNonce {
		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Printf("The mined block's hash is: \n%x", hash)
	fmt.Print("\n")

	return nonce, hash[:]
}

// Validate validates block's PoW
func (pow *proofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
