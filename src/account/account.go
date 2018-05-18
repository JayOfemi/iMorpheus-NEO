package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
	"fmt"
)

const version = byte(0x00)
const addressChecksumLen = 4

// Account stores private and public keys
type Account struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewAccount creates and returns a Account
func NewAccount() *Account {
	private, public := newKeyPair()
	account := Account{private, public}

	return &account
}

// GetAddress returns account address
func (w Account) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

// HashPubKey hashes public key
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// ValidateAddress check if address if valid
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// Checksum generates a checksum for a public key
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	///test---------------------------------------------------------------------------------------
	fmt.Println("The public key is:")
	fmt.Printf("(%d, %d)\n\n", private.PublicKey.X, private.PublicKey.Y)
	fmt.Println("The private key is:")
	fmt.Printf("%s\n\n", private.D.String())

	//words := PrivKey2Words(private.D)
	//fmt.Println("The private key words is:")
	//fmt.Println(words)
	//fmt.Printf("\n")
	//k := Words2PrivKey(words)
	//fmt.Println("The private key is:")
	//fmt.Printf("%v\n\n", k)
	///end test------------------------------------------------------------------------------------

	return *private, pubKey
}
