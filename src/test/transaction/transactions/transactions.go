package transactions

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"blockChainMorp/accounts"
	"blockChainMorp/utxo"
	"blockChainMorp/transaction"
	"blockChainMorp/types"
	"blockChainMorp/merkle_tree"
	"blockChainMorp/blockchain"
	"blockChainMorp/encrypt/secp256k1"
)

const subsidy = 50

//// Transaction represents a coin transaction
//type Transaction struct {
//	ID   []byte
//	Vin  []TXInput
//	Vout []TXOutput
//}

// TXInput represents a transaction input
type trans struct {
}

// NewTrans create a new transactions
func NewTrans() *trans {
	txs := new(trans)
	return txs
}

// IsCoinbase checks whether the transaction is coinbase
func (t trans) IsCoinbase(tx types.Transaction) bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// Serialize returns a serialized Transaction
func (t trans) Serialize(tx types.Transaction) []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

// Hash returns the hash of the Transaction
func (t *trans) Hash(tx *types.Transaction) []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(t.Serialize(txCopy))

	return hash[:]
}

// Sign signs each input of a Transaction
func (t *trans) Sign(tx *types.Transaction, privKey ecdsa.PrivateKey, prevTXs map[string]types.Transaction) {
	if t.IsCoinbase(*tx) {
		return
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := t.TrimmedCopy(tx)

	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
		txCopy.Vin[inID].PubKey = nil
	}
}

// TrimmedCopy creates a trimmed copy of Transaction to be used in signing
func (t *trans) TrimmedCopy(tx *types.Transaction) types.Transaction {
	var inputs []types.TXInput
	var outputs []types.TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, types.TXInput{vin.Txid, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, types.TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := types.Transaction{tx.ID, inputs, outputs}

	return txCopy
}

// Verify verifies signatures of Transaction inputs
func (t *trans) Verify(tx *types.Transaction, prevTXs map[string]types.Transaction, useBtc bool) bool {
	if t.IsCoinbase(*tx) {
		return true
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := t.TrimmedCopy(tx)
	curve := elliptic.P256()
	if useBtc {
		curve = secp256k1.S256()
	}

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
		txCopy.Vin[inID].PubKey = nil
	}

	return true
}

// NewCoinbaseTX creates a new coinbase transaction
func (t trans) NewCoinbaseTX(out transaction.ITXOutput, to, data string) *types.Transaction {
	if data == "" {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}

	txin := types.TXInput{[]byte{}, -1, nil, []byte(data)}
	txout := out.NewTXOutput(subsidy, to)
	tx := types.Transaction{nil, []types.TXInput{txin}, []types.TXOutput{*txout}}
	tx.ID = t.Hash(&tx)

	return &tx
}

// NewUTXOTransaction creates a new transaction
func (t *trans) NewUTXOTransaction(utxo utxo.IUTXO, out transaction.ITXOutput, account *accounts.Account, to string, amount int, UTXOSet *types.UTXOSet) *types.Transaction {
	var inputs []types.TXInput
	var outputs []types.TXOutput

	pubKeyHash := accounts.HashPubKey(account.PublicKey)
	acc, validOutputs := utxo.FindSpendableOutputs(*UTXOSet, pubKeyHash, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := types.TXInput{txID, out, nil, account.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	from := fmt.Sprintf("%s", account.GetAddress())
	outputs = append(outputs, *out.NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *out.NewTXOutput(acc-amount, from)) // a change
	}

	tx := types.Transaction{nil, inputs, outputs}
	tx.ID = t.Hash(&tx)
	blockchain.SignTransaction(UTXOSet.Blockchain, t, &tx, account.PrivateKey)

	return &tx
}

// DeserializeTransaction deserializes a transaction
func (t trans) DeserializeTransaction(data []byte) types.Transaction {
	var transaction types.Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}

	return transaction
}

// HashTransactions returns a hash of the transactions in the block
func (t trans) HashTransactions(b *types.Block) []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, t.Serialize(*tx))
	}
	mTree := merkle_tree.NewMerkleTree(transactions)

	return mTree.RootNode.Data
}
