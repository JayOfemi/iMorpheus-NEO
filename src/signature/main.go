// this is an example to show you how to check a transaction script
// and decode it.

package main

import (
	"fmt"
	"log"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

func main() {
	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8332",
		User:         "rpc",
		Pass:         "rpc",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)

	}
	defer client.Shutdown()

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	blockHash, err := client.GetBlockHash(525375)
	fmt.Println(blockHash)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)
	block, err := client.GetBlock(blockHash)
	fmt.Println(block.BlockHash())
	TxHash, err := block.TxHashes()
	//transaction count and it's hash in a block
	for index, hash := range TxHash {
		fmt.Println(index, "----", hash)
	}
	//parse transaction
	for i := 0; i < len(block.Transactions); i++ {
		if i == 268 {
			fmt.Println(block.Transactions[i].TxHash())
			//fmt.Println(block.Transactions[i].WitnessHash())
			var temp *wire.TxIn
			temp = block.Transactions[i].TxIn[0]
			//hash and index of an utxo
			fmt.Println(temp.PreviousOutPoint.Hash, "----", temp.PreviousOutPoint.Index)
			fmt.Println(temp.SignatureScript, " ", temp.Sequence)
			rawscript, err := client.DecodeScript(temp.SignatureScript)
			//details from json_rpc_api.md +231 rows
			//https://blockchain.info/tx/b1046a1bb0d168bd0b4057bb10fbb0feec09f84a89684c379a7e22bcd81c91d1
			if err == nil {
				fmt.Println(rawscript.Addresses)
				fmt.Println(rawscript.Asm)
				fmt.Println(rawscript.P2sh)
				fmt.Println(rawscript.ReqSigs)
				fmt.Println(rawscript.Type)
			}
			var tempout *wire.TxOut
			tempout = block.Transactions[i].TxOut[0]
			rawscripout, err := client.DecodeScript(tempout.PkScript)
			if err == nil {
				fmt.Println(rawscripout.Addresses) //wallet address
				fmt.Println(rawscripout.Asm)       //(pubk)-->sha256-->sha160
				fmt.Println(rawscripout.P2sh)      //hash of script
				fmt.Println(rawscripout.ReqSigs)   //count of demanding signature
				fmt.Println(rawscripout.Type)      //script type pubkeyhash
			}
		}
	}

}
