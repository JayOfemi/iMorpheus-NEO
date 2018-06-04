package client

import (
	"blockChainMorp/accounts"
	"blockChainMorp/blockchain"
	"blockChainMorp/encrypt/base58"
	"blockChainMorp/pow/proofofwork"
	"blockChainMorp/server"
	"blockChainMorp/transaction/transactions"
	"blockChainMorp/types"
	"blockChainMorp/utxo/utxo_set"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var Bc *types.Blockchain

// CreateAccount creates a new account
func CreateAccount() {
	accountsVar, _ := accounts.NewAccounts()
	address := accountsVar.CreateAccount()
	accountsVar.SaveToFile()

	fmt.Printf("Your new address: %s\n", address)

	CreateBlockchain(address)
}

// CreateBitcoinAccount creates a new account
func CreateBitcoinAccount() {
	accountsVar, _ := accounts.NewBTCAccounts()
	address := accountsVar.CreateBTCAccount()
	accountsVar.SaveBTCKeyToFile()

	fmt.Printf("Your new Bitcoin address: %s\n", address)

	CreateBlockchain(address)
}

func ListAddresses() {
	accountsVar, err := accounts.NewAccounts()
	if err != nil {
		log.Panic(err)
	}
	addresses := accountsVar.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func ListBTCAddresses() {
	accountsVar, err := accounts.NewBTCAccounts()
	if err != nil {
		log.Panic(err)
	}
	addresses := accountsVar.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func GetBalance(address string) {
	if !accounts.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := blockchain.NewBlockchain()
	UTXOSet := types.UTXOSet{bc}
	defer bc.Db.Close()

	utxos := utxoset.NewUTXOSet()
	base58coder := base58.NewBase58Coder()
	balance := 0
	pubKeyHash := base58coder.Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := utxos.FindUTXO(UTXOSet, pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func CreateBlockchain(address string) {
	tx := transactions.NewTrans()
	out := transactions.NewTXOutput()
	utxos := utxoset.NewUTXOSet()
	powVar := proofofwork.NewProofOfWork()

	bc := blockchain.CreateBlockchain(out, powVar, tx, address)
	defer bc.Db.Close()

	UTXOSet := types.UTXOSet{bc}
	utxos.Reindex(UTXOSet)

	fmt.Println("Done!")
}

func PrintChain() {
	bc := blockchain.NewBlockchain()
	defer bc.Db.Close()

	bci := blockchain.Iterator(bc)

	for {
		block := bci.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)

		powVar := proofofwork.NewProofOfWork()
		powVar.InitBlock(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(powVar.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func Send(from, to string, amount int, mineNow bool, useBtc bool) {
	if !accounts.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !accounts.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := blockchain.NewBlockchain()
	var UTXOSet types.UTXOSet
	if bc != nil {
		UTXOSet = types.UTXOSet{bc}
	} else {
		return
	}

	defer bc.Db.Close()

	var accountsVar *accounts.Accounts
	var err error
	if useBtc {
		accountsVar, err = accounts.NewBTCAccounts()
		if err != nil {
			log.Panic(err)
		}
	} else {
		accountsVar, err = accounts.NewAccounts()
		if err != nil {
			log.Panic(err)
		}
	}
	acount := accountsVar.GetAccount(from)

	itx := transactions.NewTrans()
	out := transactions.NewTXOutput()
	utxo := utxoset.NewUTXOSet()
	tx := itx.NewUTXOTransaction(utxo, out, &acount, to, amount, &UTXOSet)

	if mineNow {
		powVar := proofofwork.NewProofOfWork()

		cbTx := itx.NewCoinbaseTX(out, from, "")
		txs := []*types.Transaction{cbTx, tx}

		newBlock := blockchain.MineBlock(powVar, bc, itx, txs, useBtc)
		utxo.Update(UTXOSet, newBlock)

		fmt.Println("Sent and mined successfully!")
	} else {
		server.SendTx(itx, server.KnownNodes[0], tx)
		fmt.Println("Sent a transaction.")
	}
}

func SendAndDo(from, to string, amount int, mineNow bool, useBtc bool, dof func()) {
	if !accounts.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !accounts.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := blockchain.NewBlockchain()
	var UTXOSet types.UTXOSet
	if bc != nil {
		UTXOSet = types.UTXOSet{bc}
	} else {
		return
	}

	defer func() {
		if bc != nil {
			bc.Db.Close()
			blockchain.DbFileOpenedFlag = false

			dof()
		}
	}()

	var accountsVar *accounts.Accounts
	var err error
	if useBtc {
		accountsVar, err = accounts.NewBTCAccounts()
		if err != nil {
			log.Panic(err)
		}
	} else {
		accountsVar, err = accounts.NewAccounts()
		if err != nil {
			log.Panic(err)
		}
	}
	acount := accountsVar.GetAccount(from)

	itx := transactions.NewTrans()
	out := transactions.NewTXOutput()
	utxo := utxoset.NewUTXOSet()
	tx := itx.NewUTXOTransaction(utxo, out, &acount, to, amount, &UTXOSet)

	if mineNow {
		powVar := proofofwork.NewProofOfWork()

		cbTx := itx.NewCoinbaseTX(out, from, "")
		txs := []*types.Transaction{cbTx, tx}

		newBlock := blockchain.MineBlock(powVar, bc, itx, txs, useBtc)
		utxo.Update(UTXOSet, newBlock)

		fmt.Println("Sent and mined successfully!")
	} else {
		server.SendTx(itx, server.KnownNodes[0], tx)
		fmt.Println("Sent a transaction.")
	}
}

func StartNode(minerAddress string, useBtc bool) {
	fmt.Printf("Starting node\n")

	if len(minerAddress) > 0 {
		if (strings.Compare(minerAddress, "") == 0) ||
			(len(minerAddress) <= accounts.AddressChecksumLen*2) || (!accounts.ValidateAddress(minerAddress)) {
			log.Panic("Wrong miner address!")
		} else {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		}
	}

	txVar := transactions.NewTrans()
	outVar := transactions.NewTXOutput()
	utxoVar := utxoset.NewUTXOSet()
	powVar := proofofwork.NewProofOfWork()

	server.StartServer(txVar, outVar, powVar, utxoVar, minerAddress, Bc, useBtc)
}

func ReindexUTXO() {
	utxoVar := utxoset.NewUTXOSet()
	bc := blockchain.NewBlockchain()
	UTXOSet := types.UTXOSet{bc}
	utxoVar.Reindex(UTXOSet)

	count := utxoVar.CountTransactions(UTXOSet)
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}
