package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const accountFile = "account.dat"

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

// CreateAccount adds a Account to Accounts
func (ws *Accounts) CreateAccount() string {
	account := NewAccount()
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
	accountFile := fmt.Sprintf(accountFile)
	if _, err := os.Stat(accountFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(accountFile)
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

// SaveToFile saves accounts to a file
func (ws Accounts) SaveToFile() {
	var content bytes.Buffer
	accountFile := fmt.Sprintf(accountFile)

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(accountFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
