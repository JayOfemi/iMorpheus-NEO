package transactions

import (
	"bytes"
	"blockChainMorp/accounts"
	"blockChainMorp/types"
)

//// TXInput represents a transaction input
//type TXInput struct {
//	Txid      []byte
//	Vout      int
//	Signature []byte
//	PubKey    []byte
//}

type txInput struct {
}

// NewTXInput create a new transaction input
func NewTXInput() *txInput {
	txInput := new(txInput)
	return txInput
}

// UsesKey checks whether the address initiated the transaction
func (in *txInput) UsesKey(input types.TXInput, pubKeyHash []byte) bool {
	lockingHash := accounts.HashPubKey(input.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
