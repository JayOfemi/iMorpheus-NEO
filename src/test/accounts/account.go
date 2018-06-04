package accounts

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
	"blockChainMorp/encrypt/base58"
	"blockChainMorp/encrypt/secp256k1"
)

const version = byte(0x00)
const AddressChecksumLen = 4

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
	checksumVar := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksumVar...)

	base58coder := base58.NewBase58Coder()
	address := base58coder.Encode(fullPayload)

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
	base58coder := base58.NewBase58Coder()
	pubKeyHash := base58coder.Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-AddressChecksumLen:]
	versionVar := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-AddressChecksumLen]
	targetChecksum := checksum(append([]byte{versionVar}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// Checksum generates a checksum for a public key
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:AddressChecksumLen]
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

// NewBTCAccount creates and returns a Account of BTC
func NewBTCAccount() *Account {
	private, public := newBTCKeyPair()
	account := Account{private, public}

	return &account
}

func newBTCKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := secp256k1.S256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}
