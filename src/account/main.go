package main

import "fmt"

func main() {
	accounts, _ := NewAccounts()
	address := accounts.CreateAccount()
	accounts.SaveToFile()

	fmt.Printf("The new address is:\n%s\n", address)
}
