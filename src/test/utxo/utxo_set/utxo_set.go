package utxoset

import (
	"encoding/hex"
	"log"

	"github.com/boltdb/bolt"
	"blockChainMorp/blockchain"
	"blockChainMorp/types"
	"blockChainMorp/transaction/transactions"
)

const utxoBucket = "chainstate"

// UTXOSet represents UTXO set
type utxoSet struct {
}

// NewUTXOSet create a new utxo set
func NewUTXOSet() *utxoSet {
	utxos := new(utxoSet)
	return utxos
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (u utxoSet) FindSpendableOutputs(utxo types.UTXOSet, pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := utxo.Blockchain.Db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		iouts := transactions.NewTXOutputs()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := iouts.DeserializeOutputs(v)

			iout := transactions.NewTXOutput()
			for outIdx, out := range outs.Outputs {
				if iout.IsLockedWithKey(out, pubkeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return accumulated, unspentOutputs
}

// FindUTXO finds UTXO for a public key hash
func (u utxoSet) FindUTXO(set types.UTXOSet, pubKeyHash []byte) []types.TXOutput {
	var UTXOs []types.TXOutput
	db := set.Blockchain.Db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		iouts := transactions.NewTXOutputs()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := iouts.DeserializeOutputs(v)

			iout := transactions.NewTXOutput()
			for _, out := range outs.Outputs {
				if iout.IsLockedWithKey(out, pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// CountTransactions returns the number of transactions in the UTXO set
func (u utxoSet) CountTransactions(set types.UTXOSet) int {
	db := set.Blockchain.Db
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return counter
}

// Reindex rebuilds the UTXO set
func (u utxoSet) Reindex(set types.UTXOSet) {
	db := set.Blockchain.Db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}

		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	itx := transactions.NewTrans()
	UTXO := blockchain.FindUTXO(set.Blockchain, itx)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}

			iouts := transactions.NewTXOutputs()
			err = b.Put(key, iouts.Serialize(outs))
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
}

// Update updates the UTXO set with transactions from the Block
// The Block is considered to be the tip of a blockchain
func (u utxoSet) Update(set types.UTXOSet, block *types.Block) {
	db := set.Blockchain.Db

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		itx := transactions.NewTrans()
		iouts := transactions.NewTXOutputs()
		for _, tx := range block.Transactions {
			if itx.IsCoinbase(*tx) == false {
				for _, vin := range tx.Vin {
					updatedOuts := types.TXOutputs{}
					outsBytes := b.Get(vin.Txid)
					outs := iouts.DeserializeOutputs(outsBytes)

					for outIdx, out := range outs.Outputs {
						if outIdx != vin.Vout {
							updatedOuts.Outputs = append(updatedOuts.Outputs, out)
						}
					}

					if len(updatedOuts.Outputs) == 0 {
						err := b.Delete(vin.Txid)
						if err != nil {
							log.Panic(err)
						}
					} else {
						err := b.Put(vin.Txid, iouts.Serialize(updatedOuts))
						if err != nil {
							log.Panic(err)
						}
					}

				}
			}

			newOutputs := types.TXOutputs{}
			for _, out := range tx.Vout {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}

			err := b.Put(tx.ID, iouts.Serialize(newOutputs))
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
