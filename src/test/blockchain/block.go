package blockchain

import (
	"blockChainMorp/pow"
	"blockChainMorp/types"
	"bytes"
	"encoding/gob"
	"log"
	"time"
	"fmt"
)

// NewBlock creates and returns Block
func NewBlock(powVar pow.IProofOfWork, transactions []*types.Transaction, prevBlockHash []byte, height int) *types.Block {
	block := &types.Block {
		time.Now().Unix(),
		transactions,
		prevBlockHash,
		[]byte{},
		0,
		height,
	}

	powVar.InitBlock(block)
	nonce, hash := powVar.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(powVar pow.IProofOfWork, coinbase *types.Transaction) *types.Block {
	genesisBlock := NewBlock(powVar, []*types.Transaction{coinbase}, []byte{}, 0)
	fmt.Println("A new genesis block created!")
	return genesisBlock
}

// SerializeBlock serializes the block
func SerializeBlock(b *types.Block) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock deserializes a block
func DeserializeBlock(d []byte) *types.Block {
	var block types.Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
