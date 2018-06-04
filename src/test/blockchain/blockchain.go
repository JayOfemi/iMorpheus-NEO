package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"blockChainMorp/pow"
	"blockChainMorp/transaction"
	"blockChainMorp/types"
	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

var DbFileOpenedFlag = false

// CreateBlockchain creates a new blockchain DB
func CreateBlockchain(out transaction.ITXOutput, powVar pow.IProofOfWork, itx transaction.ITransaction, address string) *types.Blockchain {
	if dbExists(dbFile) {
		fmt.Println("Blockchain already exists, Exit.")
		os.Exit(0)
	}

	var tip []byte

	cbtx := itx.NewCoinbaseTX(out, address, genesisCoinbaseData)
	genesis := NewGenesisBlock(powVar, cbtx)

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, SerializeBlock(genesis))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := types.Blockchain{tip, db}

	return &bc
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain() *types.Blockchain {
	if dbExists(dbFile) == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	if !DbFileOpenedFlag {
		var tip []byte
		db, err := bolt.Open(dbFile, 0600, nil)
		if err != nil {
			log.Panic(err)
		}

		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			tip = b.Get([]byte("l"))

			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		bc := types.Blockchain{tip, db}
		DbFileOpenedFlag = true

		return &bc
	} else {
		fmt.Println("The block chain data file has been opened, do not opened again!")
	}

	return nil
}

// AddBlock saves the block into the blockchain
func AddBlock(bc *types.Blockchain, block *types.Block) {
	err := bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		blockInDb := b.Get(block.Hash)

		if blockInDb != nil {
			return nil
		}

		blockData := SerializeBlock(block)
		err := b.Put(block.Hash, blockData)
		if err != nil {
			log.Panic(err)
		}

		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)

		if block.Height > lastBlock.Height {
			err = b.Put([]byte("l"), block.Hash)
			if err != nil {
				log.Panic(err)
			}
			bc.Tip = block.Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// FindTransaction finds a transaction by its ID
func FindTransaction(bc *types.Blockchain, ID []byte) (types.Transaction, error) {
	bci := Iterator(bc)

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return types.Transaction{}, errors.New("Transaction is not found")
}

// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
func FindUTXO(bc *types.Blockchain, itx transaction.ITransaction) map[string]types.TXOutputs {
	UTXO := make(map[string]types.TXOutputs)
	spentTXOs := make(map[string][]int)
	bci := Iterator(bc)

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if itx.IsCoinbase(*tx) == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// Iterator returns a BlockchainIterat
func Iterator(bc *types.Blockchain) *BlockchainIterator {
	bci := &BlockchainIterator{bc.Tip, bc.Db}

	return bci
}

// GetBestHeight returns the height of the latest block
func GetBestHeight(bc *types.Blockchain) int {
	var lastBlock types.Block

	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Height
}

// GetBlock finds a block by its hash and returns it
func GetBlock(bc *types.Blockchain, blockHash []byte) (types.Block, error) {
	var block types.Block

	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		blockData := b.Get(blockHash)

		if blockData == nil {
			return errors.New("Block is not found.")
		}

		block = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
}

// GetBlockHashes returns a list of hashes of all the blocks in the chain
func GetBlockHashes(bc *types.Blockchain) [][]byte {
	var blocks [][]byte
	bci := Iterator(bc)

	for {
		block := bci.Next()

		blocks = append(blocks, block.Hash)

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return blocks
}

// MineBlock mines a new block with the provided transactions
func MineBlock(powVar pow.IProofOfWork, bc *types.Blockchain, itx transaction.ITransaction,
	transactions []*types.Transaction, useBtc bool) *types.Block {
	var lastHash []byte
	var lastHeight int

	for _, tx := range transactions {
		// TODO: ignore transaction if it's not valid
		if VerifyTransaction(bc, itx, tx, useBtc) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		blockData := b.Get(lastHash)
		block := DeserializeBlock(blockData)

		lastHeight = block.Height

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(powVar, transactions, lastHash, lastHeight+1)

	err = bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, SerializeBlock(newBlock))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.Tip = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return newBlock
}

// SignTransaction signs inputs of a Transaction
func SignTransaction(bc *types.Blockchain, itx transaction.ITransaction, tx *types.Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]types.Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := FindTransaction(bc, vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	itx.Sign(tx, privKey, prevTXs)
}

// VerifyTransaction verifies transaction input signatures
func VerifyTransaction(bc *types.Blockchain, itx transaction.ITransaction, tx *types.Transaction, useBtc bool) bool {
	if itx.IsCoinbase(*tx) {
		return true
	}

	prevTXs := make(map[string]types.Transaction)

	var prevTXOut = 0
	for _, vin := range tx.Vin {
		prevTX, err := FindTransaction(bc, vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX

		for id, out := range prevTX.Vout {
			if id == vin.Vout {
				prevTXOut += out.Value
			}
		}
	}

	var currtTXOut = 0
	for _, out := range tx.Vout {
		currtTXOut += out.Value
	}

	if prevTXOut == currtTXOut {
		return itx.Verify(tx, prevTXs, useBtc)
	} else {
		return false
	}
}

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
