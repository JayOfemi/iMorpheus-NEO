package utxo

import (
	"blockChainMorp/types"
)

type IUTXO interface {

	// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
	FindSpendableOutputs(set types.UTXOSet, pubkeyHash []byte, amount int) (int, map[string][]int)

	// FindUTXO finds UTXO for a public key hash
	FindUTXO(set types.UTXOSet, pubKeyHash []byte) []types.TXOutput

	// CountTransactions returns the number of transactions in the UTXO set
	CountTransactions(set types.UTXOSet) int

	// Reindex rebuilds the UTXO set
	Reindex(set types.UTXOSet)

	// Update updates the UTXO set with transactions from the Block
	// The Block is considered to be the tip of a blockchain
	Update(set types.UTXOSet, block *types.Block)
}