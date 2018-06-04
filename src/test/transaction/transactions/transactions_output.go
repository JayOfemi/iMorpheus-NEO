package transactions

import (
	"bytes"
	"encoding/gob"
	"log"
	"blockChainMorp/encrypt/base58"
	"blockChainMorp/types"
)

//// TXOutput represents a transaction output
//type TXOutput struct {
//	Value      int
//	PubKeyHash []byte
//}

type txOutput struct {
}

// NewTXOutput create a new transaction output
func NewTXOutput() *txOutput {
	txOutput := new(txOutput)
	return txOutput
}

// Lock signs the output
func (out *txOutput) Lock(output *types.TXOutput, address []byte) {
	base58coder := base58.NewBase58Coder()
	pubKeyHash := base58coder.Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	output.PubKeyHash = pubKeyHash
}

// IsLockedWithKey checks if the output can be used by the owner of the pubkey
func (out *txOutput) IsLockedWithKey(output types.TXOutput, pubKeyHash []byte) bool {
	return bytes.Compare(output.PubKeyHash, pubKeyHash) == 0
}

// NewTXOutput create a new TXOutput
func (out *txOutput) NewTXOutput(value int, address string) *types.TXOutput {
	txo := &types.TXOutput{value, nil}
	out.Lock(txo, []byte(address))

	return txo
}

//// TXOutputs collects TXOutput
//type TXOutputs struct {
//	Outputs []transaction.TXOutput
//}

type txOutputs struct {
}

// NewTXOutputs create a new transaction output
func NewTXOutputs() *txOutputs {
	txOutputs := new(txOutputs)
	return txOutputs
}

// Serialize serializes TXOutputs
func (outs txOutputs) Serialize(outputs types.TXOutputs) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outputs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// DeserializeOutputs deserializes TXOutputs
func (outs txOutputs) DeserializeOutputs(data []byte) types.TXOutputs {
	var outputs types.TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
