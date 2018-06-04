package accounts

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"blockChainMorp/encrypt/secp256k1"
)

const AccountFile = "accounts.dat"
const AccountBtcFile = "accountsBtc.dat"

// Accounts stores a collection of accounts
type Accounts struct {
	Accounts map[string]*Account
}

// NewAccounts creates Accounts and fills it from a file if it exists
func NewAccounts() (*Accounts, error) {
	accounts := Accounts{}
	accounts.Accounts = make(map[string]*Account)

	err := accounts.LoadFromFile()

	return &accounts, err
}

// NewBTCAccounts creates Accounts and fills it from a file if it exists
func NewBTCAccounts() (*Accounts, error) {
	accounts := Accounts{}
	accounts.Accounts = make(map[string]*Account)

	err := accounts.LoadBTCKeyFromFile()

	return &accounts, err
}

// CreateAccount adds a Account to Accounts
func (ws *Accounts) CreateAccount() string {
	account := NewAccount()
	address := fmt.Sprintf("%s", account.GetAddress())

	ws.Accounts[address] = account

	return address
}

// CreateBTCAccount adds a Account of BTC to Accounts
func (ws *Accounts) CreateBTCAccount() string {
	account := NewBTCAccount()
	address := fmt.Sprintf("%s", account.GetAddress())

	ws.Accounts[address] = account

	return address
}

// GetAddresses returns an array of addresses stored in the account file
func (ws *Accounts) GetAddresses() []string {
	var addresses []string

	for address := range ws.Accounts {
		addresses = append(addresses, address)
	}

	return addresses
}

// GetAccount returns a Account by its address
func (ws Accounts) GetAccount(address string) Account {
	return *ws.Accounts[address]
}

// LoadFromFile loads accounts from the file
func (ws *Accounts) LoadFromFile() error {
	AccountFile := fmt.Sprintf(AccountFile)
	if _, err := os.Stat(AccountFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(AccountFile)
	if err != nil {
		log.Panic(err)
	}

	var accounts Accounts
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&accounts)
	if err != nil {
		log.Panic(err)
	}

	ws.Accounts = accounts.Accounts

	return nil
}

// LoadBTCKeyFromFile loads accounts from the file
func (ws *Accounts) LoadBTCKeyFromFile() error {
	AccountFile := fmt.Sprintf(AccountBtcFile)
	if _, err := os.Stat(AccountFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(AccountFile)
	if err != nil {
		log.Panic(err)
	}

	var accounts Accounts
	gob.Register(secp256k1.S256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&accounts)
	if err != nil {
		log.Panic(err)
	}

	ws.Accounts = accounts.Accounts

	return nil
}

// SaveToFile saves accounts to a file
func (ws Accounts) SaveToFile() {
	var content bytes.Buffer
	AccountFile := fmt.Sprintf(AccountFile)

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(AccountFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

// SaveBTCKeyToFile saves accounts of BTC to a file
func (ws Accounts) SaveBTCKeyToFile() {
	var content bytes.Buffer
	AccountFile := fmt.Sprintf(AccountBtcFile)

	gob.Register(secp256k1.S256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(AccountFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
