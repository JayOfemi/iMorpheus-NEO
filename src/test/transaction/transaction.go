package transaction

import (
	"crypto/ecdsa"
	"blockChainMorp/accounts"
	"blockChainMorp/types"
	"blockChainMorp/utxo"
)

//ITransaction is the interface of Transaction
type ITransaction interface {

	// IsCoinbase checks whether the transaction is coinbase
	IsCoinbase(tx types.Transaction) bool

	// Serialize returns a serialized Transaction
	Serialize(tx types.Transaction) []byte

	// Hash returns the hash of the Transaction
	Hash(tx *types.Transaction) []byte

	// Sign signs each input of a Transaction
	Sign(tx *types.Transaction, privKey ecdsa.PrivateKey, prevTXs map[string]types.Transaction)

	// TrimmedCopy creates a trimmed copy of Transaction to be used in signing
	TrimmedCopy(tx *types.Transaction) types.Transaction

	// Verify verifies signatures of Transaction inputs
	Verify(tx *types.Transaction, prevTXs map[string]types.Transaction, useBtc bool) bool

	// NewCoinbaseTX creates a new coinbase transaction
	NewCoinbaseTX(out ITXOutput, to, data string) *types.Transaction

	// NewUTXOTransaction creates a new transaction
	NewUTXOTransaction(utxo utxo.IUTXO, out ITXOutput, account *accounts.Account, to string, amount int, UTXOSet *types.UTXOSet) *types.Transaction

	// DeserializeTransaction deserializes a transaction
	DeserializeTransaction(data []byte) types.Transaction

	// HashTransactions returns a hash of the transactions in the block
	HashTransactions(b *types.Block) []byte
}

type ITXInput interface {

	// UsesKey checks whether the address initiated the transaction
	UsesKey(input types.TXInput, pubKeyHash []byte) bool

}

type ITXOutput interface {

	// Lock signs the output
	Lock(output *types.TXOutput, address []byte)

	// IsLockedWithKey checks if the output can be used by the owner of the pubkey
	IsLockedWithKey(output types.TXOutput, pubKeyHash []byte) bool

	// NewTXOutput create a new TXOutput
	NewTXOutput(value int, address string) *types.TXOutput

}

type ITXOutputs interface {

	// Serialize serializes TXOutputs
	Serialize(outputs types.TXOutputs) []byte

	// DeserializeOutputs deserializes TXOutputs
	DeserializeOutputs(data []byte) types.TXOutputs
}
