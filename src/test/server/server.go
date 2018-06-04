package server

import (
	"blockChainMorp/blockchain"
	"blockChainMorp/pow"
	"blockChainMorp/transaction"
	"blockChainMorp/types"
	"blockChainMorp/utxo"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12
const minMineCnt = 1

var nodeAddress string
var miningAddress string
var KnownNodes = []string{"192.168.2.214:10308"}
var blocksInTransit [][]byte
var mempool = make(map[string]types.Transaction)
var stopFlag = true
var chs = make([]chan int, 1)
var chanNum = 0

type addr struct {
	AddrList []string
}

type block struct {
	AddrFrom string
	Block    []byte
}

type getblocks struct {
	AddrFrom string
}

type getdata struct {
	AddrFrom string
	Type     string
	ID       []byte
}

type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

type tx struct {
	AddFrom     string
	Transaction []byte
}

type verzion struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

func extractCommand(request []byte) []byte {
	return request[:commandLength]
}

func requestBlocks() {
	for _, node := range KnownNodes {
		sendGetBlocks(node)
	}
}

func sendAddr(address string) {
	nodes := addr{KnownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := gobEncode(nodes)
	request := append(commandToBytes("addr"), payload...)

	sendData(address, request)
}

func sendBlock(addr string, b *types.Block) {
	data := block{nodeAddress, blockchain.SerializeBlock(b)}
	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)

	sendData(addr, request)
}

func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available and will be removed!\n", addr)
		var updatedNodes []string

		//remove this node
		for _, node := range KnownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		KnownNodes = updatedNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func sendInv(address, kind string, items [][]byte) {
	inventory := inv{nodeAddress, kind, items}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

func sendGetBlocks(address string) {
	payload := gobEncode(getblocks{nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func sendGetData(address, kind string, id []byte) {
	payload := gobEncode(getdata{nodeAddress, kind, id})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}

func sendTx(itx transaction.ITransaction, addr string, tnx *types.Transaction) {
	data := tx{nodeAddress, itx.Serialize(*tnx)}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	sendData(addr, request)
}

func sendVersion(addr string, bc *types.Blockchain) {
	bestHeight := blockchain.GetBestHeight(bc)
	payload := gobEncode(verzion{nodeVersion, bestHeight, nodeAddress})

	request := append(commandToBytes("version"), payload...)

	sendData(addr, request)
}

func handleAddr(request []byte) {
	var buff bytes.Buffer
	var payload addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	KnownNodes = append(KnownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(KnownNodes))
	requestBlocks()
}

func handleBlock(iutxo utxo.IUTXO, request []byte, bc *types.Blockchain) {
	var buff bytes.Buffer
	var payload block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := blockchain.DeserializeBlock(blockData)

	fmt.Println("Recevied a new block!")
	blockchain.AddBlock(bc, block)

	fmt.Printf("Added block %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := types.UTXOSet{bc}
		iutxo.Reindex(UTXOSet)
	}
}

func handleInv(request []byte) {
	var buff bytes.Buffer
	var payload inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		var newInTransit [][]byte
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func handleGetBlocks(request []byte, bc *types.Blockchain) {
	var buff bytes.Buffer
	var payload getblocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := blockchain.GetBlockHashes(bc)
	sendInv(payload.AddrFrom, "block", blocks)
}

func handleGetData(itx transaction.ITransaction, request []byte, bc *types.Blockchain) {
	var buff bytes.Buffer
	var payload getdata

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := blockchain.GetBlock(bc, []byte(payload.ID))
		if err != nil {
			return
		}

		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		sendTx(itx, payload.AddrFrom, &tx)
		// delete(mempool, txID)
	}
}

func handleTx(itx transaction.ITransaction, iout transaction.ITXOutput, ipow pow.IProofOfWork,
	iutxo utxo.IUTXO, request []byte, bc *types.Blockchain, useBtc bool) {
	var buff bytes.Buffer
	var payload tx

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := itx.DeserializeTransaction(txData)
	mempool[hex.EncodeToString(tx.ID)] = tx

	if nodeAddress == KnownNodes[0] {
		for _, node := range KnownNodes {
			if node != nodeAddress && node != payload.AddFrom {
				sendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(mempool) >= minMineCnt && len(miningAddress) > 0 {
		MineTransactions:
			var txs []*types.Transaction

			for id := range mempool {
				tx := mempool[id]
				if blockchain.VerifyTransaction(bc, itx, &tx, useBtc) {
					txs = append(txs, &tx)
				}
			}

			if len(txs) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}

			cbTx := itx.NewCoinbaseTX(iout, miningAddress, "")
			txs = append(txs, cbTx)

			newBlock := blockchain.MineBlock(ipow, bc, itx, txs, useBtc)
			UTXOSet := types.UTXOSet{bc}
			iutxo.Reindex(UTXOSet)

			fmt.Println("New block is mined!")

			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(mempool, txID)
			}

			for _, node := range KnownNodes {
				if node != nodeAddress {
					sendInv(node, "block", [][]byte{newBlock.Hash})
				}
			}

			if len(mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}

func handleVersion(request []byte, bc *types.Blockchain) {
	var buff bytes.Buffer
	var payload verzion

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := blockchain.GetBestHeight(bc)
	foreignerBestHeight := payload.BestHeight

	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	// sendAddr(payload.AddrFrom)
	if !nodeIsKnown(payload.AddrFrom) {
		KnownNodes = append(KnownNodes, payload.AddrFrom)
	}
}

func handleConnection(itx transaction.ITransaction, iout transaction.ITXOutput, ipow pow.IProofOfWork,
	iutxo utxo.IUTXO, conn net.Conn, bc *types.Blockchain, useBtc bool, chNum int) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "addr":
		handleAddr(request)
	case "block":
		handleBlock(iutxo, request, bc)
	case "inv":
		handleInv(request)
	case "getblocks":
		handleGetBlocks(request, bc)
	case "getdata":
		handleGetData(itx, request, bc)
	case "tx":
		handleTx(itx, iout, ipow, iutxo, request, bc, useBtc)
	case "version":
		handleVersion(request, bc)
	case "stopserver":

	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()

	chs[chNum] <- 1
}

func SendTx(itx transaction.ITransaction, addr string, tnx *types.Transaction) {
	sendTx(itx, addr, tnx)
}

// StartServer starts a node
func StartServer(itx transaction.ITransaction, iout transaction.ITXOutput, ipow pow.IProofOfWork,
	iutxo utxo.IUTXO, minerAddress string, bc *types.Blockchain, useBtc bool) {
	if stopFlag == false {
		return
	}
	stopFlag = false
	nodeAddress = fmt.Sprintf("%s:10308", getInternalIP())
	miningAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if bc != nil {
			bc.Db.Close()
			blockchain.DbFileOpenedFlag = false
		}
		ln.Close()
	}()

	bc = blockchain.NewBlockchain()

	if nodeAddress != KnownNodes[0] {
		sendVersion(KnownNodes[0], bc)
	}

	chanNum = 0
	if chs[0] == nil {
		chs[0] = make(chan int)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		} else {
			go handleConnection(itx, iout, ipow, iutxo, conn, bc, useBtc, chanNum)
			chanNum++
		}

		if stopFlag {
			return
		}

		chs = append(chs, make(chan int))
	}
}

func StopServer() {
	stopFlag = true

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:10308", getInternalIP()))
	if err != nil {
		log.Panic(err)
	}

	// stop ln.Accept()
	_, err = conn.Write([]byte("stopserver"))
	if err != nil {
		log.Panic(err)
	}
	conn.Close()

	for _, ch := range chs {
		<-ch
	}

	fmt.Println("Server Stopped!")
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func nodeIsKnown(addr string) bool {
	for _, node := range KnownNodes {
		if node == addr {
			return true
		}
	}

	return false
}

func getInternalIP() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		os.Stderr.WriteString("Oops:" + err.Error())
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	log.Panic("Oops: getInternalIP() failed!")
	return string("0")
}
