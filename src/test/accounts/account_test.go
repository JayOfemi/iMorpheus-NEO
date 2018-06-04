package accounts

import (
	"fmt"
	"math/big"
	"testing"
)

func TestAccount_GetAddress(t *testing.T) {
	accountsVar, _ := NewAccounts()
	address := accountsVar.CreateAccount()
	//accountsVar.SaveToFile()

	privkey := accountsVar.GetAccount(address).PrivateKey.D
	publickey := big.NewInt(0).SetBytes(accountsVar.GetAccount(address).PublicKey)
	var strOut1 = fmt.Sprintf("{PublicKey: \"%s\",\nPrivateKey: \"%s\",\nMnemonics: [\"", publickey.String(), privkey.String())

	words := PrivKey2Words(privkey)


	for i := 0; i < wordsCount; i++ {
		if i != wordsCount - 1 {
			strOut1 += fmt.Sprintf("%s\", \"", words[i])
		} else {
			strOut1 += fmt.Sprintf("%s\"],\n", words[i])
		}
	}
	strOut1 += fmt.Sprintf("Address: \"%s\"}\n", address)
	fmt.Printf(strOut1)

	k := Words2PrivKey(words)
	//fmt.Println("The mnemonics is:")
	//fmt.Println(words)
	//fmt.Printf("\n")
	//fmt.Println("The private key is:")
	//fmt.Printf("%v\n\n", k)
	//fmt.Printf("The new address is:\n%s\n", address)

	if k.Cmp(privkey) != 0 {
		t.Log("Error!")
		t.Fail()
	}
}

func TestBTCAccount_GetAddress(t *testing.T) {
	accountsVar, _ := NewBTCAccounts()
	address := accountsVar.CreateBTCAccount()
	//accountsVar.SaveBTCKeyToFile()

	privkey := accountsVar.GetAccount(address).PrivateKey.D
	publickey := big.NewInt(0).SetBytes(accountsVar.GetAccount(address).PublicKey)
	var strOut1 = fmt.Sprintf("{PublicKey: \"%s\",\nPrivateKey: \"%s\",\nMnemonics: [\"", publickey.String(), privkey.String())

	words := PrivKey2Words(privkey)


	for i := 0; i < wordsCount; i++ {
		if i != wordsCount - 1 {
			strOut1 += fmt.Sprintf("%s\", \"", words[i])
		} else {
			strOut1 += fmt.Sprintf("%s\"],\n", words[i])
		}
	}
	strOut1 += fmt.Sprintf("Address: \"%s\"}\n", address)
	fmt.Printf(strOut1)

	k := Words2PrivKey(words)
	//fmt.Println("The mnemonics is:")
	//fmt.Println(words)
	//fmt.Printf("\n")
	//fmt.Println("The private key is:")
	//fmt.Printf("%v\n\n", k)
	//fmt.Printf("The new Bitcoin address is:\n%s\n", address)

	if k.Cmp(privkey) != 0 {
		t.Log("Error!")
		t.Fail()
	}
}
