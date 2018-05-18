package main

import (
	"crypto/sha256"
	"math/big"
	"fmt"
)

const addressChecksum8bitsLen = 1 //8*addressChecksum8bitsLen bits checksum
const wordsCount = 24
const wordPerBit = 11
const dictionaryWordsCnt = 1 << wordPerBit - 1

var privKeysDictionary = [dictionaryWordsCnt]string{
	"ability", "able", "about", "above", "abroad", "absence", "absent", "accent", "accept", "accident",
	"according", "account", "ache", "achieve", "across", "action", "active", "activity", "actor", "advertisement",
	"advance", "advantage", "adventure", "advertise", "actual", "advice", "advise", "aeroplane", "affair", "affect",
	"afford", "afraid", "africa", "african", "after", "afternoon", "afterwards", "again", "against", "age",
	"aggression", "aggressive", "ago", "agree", "agreement", "agricultural", "agriculture", "ahead", "aid", "aids",
	"aim", "air", "aircraft", "airline", "airmail", "airplane", "airport", "alarm", "alive", "all",
	"allow", "almost", "alone", "along", "aloud", "already", "also", "although", "altogether", "always",
	"amaze", "ambulance", "among", "amuse", "amusement", "ancestor", "ancient", "and", "anger", "angry",
	"animal", "announce", "announcement", "annoy", "another", "answer", "ant", "antarctic", "antique", "anxious",
	"anybody", "anyhow", "anyone", "anything", "anyway", "anywhere", "apartment", "apologize", "apology", "appear",
	"appearance", "apple", "application", "apply", "appointment", "arctic", "area", "argue", "argument", "arise",
	"arithmetic", "arm", "armchair", "army", "arose", "around", "arrange", "art", "arrow", "arrive",
	"article", "artist", "aside", "ask", "asleep", "assistant", "astonish", "astronaut", "astronomy", "athlete",
	"atmosphere", "atom", "attack", "attempt", "attend", "attention", "attitude", "attract", "attractive", "audience",
	"baby", "back", "backache", "background", "backward", "bacon", "bacteria", "bacterium", "bad", "badly",
	"baggage", "bake", "bakery", "balance", "balcony", "ball", "ballet", "balloon", "ballpoint", "bamboo",
	"barbecue", "barber", "barbershop", "bargain", "bark", "base", "baseball", "basement", "basic", "basin",
	"basket", "basketball", "bat", "bath", "bathe", "bathrobe", "bathroom", "bathtub", "battery", "battle",
	"beast", "beat", "beaten", "beautiful", "beauty", "became", "because", "become", "bed", "bedroom",
	"bee", "beef", "beehive", "beer", "before", "beg", "began", "begin", "beginning", "begun",


}

// Generates a 16bits checksum for input
func checksum8(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksum8bitsLen]
}

// Verify the checksum is correct
func verifyChecksum(code []byte) bool {
	checksum := checksum8(code[addressChecksum8bitsLen:])
	if len(checksum) != addressChecksum8bitsLen {
		return false
	}

	for k, v := range checksum {
		if v != code[k] {
			return false
		}
	}

	return true
}

//Encode code to wordValue
func encode2WordValue(code []byte) [wordsCount]int64 {
	var wordValue [wordsCount]int64

	for i := 0; i < wordsCount; i++ {
		valSrc := big.NewInt(0).SetBytes(code)
		valTemp := valSrc.Rsh(valSrc, wordPerBit * uint(i))
		valAnd:= big.NewInt(dictionaryWordsCnt)
		val := valTemp.And(valTemp, valAnd)
		wordValue[i] = val.Int64()
	}

	return wordValue
}

//Decode wordValue to code
func decode2Code(wordValue [wordsCount]int64) []byte {
	var code []byte

	val := big.NewInt(0)
	for i := 0; i < wordsCount; i++ {
		valAdd := big.NewInt(wordValue[wordsCount-i-1])
		val = val.Lsh(val, wordPerBit)
		val = val.Add(val, valAdd)
	}

	code = val.Bytes()
	return code
}

//Encode private key to keywords
func PrivKey2Words(privKey *big.Int) []string {
	checksum := checksum8(privKey.Bytes())
	code := append(checksum, privKey.Bytes()...)

	wordValue := encode2WordValue(code)

	words:= make([]string, dictionaryWordsCnt)
	for i := 0; i < wordsCount; i++ {
		words[i] = privKeysDictionary[wordValue[i]]
	}

	return words
}

//Decode keywords to private key
func Words2PrivKey(words []string) *big.Int {
	wordsDictionary := make(map[string]int64, dictionaryWordsCnt)
	for i := int64(0); i < dictionaryWordsCnt; i++ {
		wordsDictionary[privKeysDictionary[i]] = i
	}

	///test-------------------------------------------------------------------------
	for j := int64(0); j < dictionaryWordsCnt; j++ {
		fmt.Println(wordsDictionary[privKeysDictionary[j]])
		if wordsDictionary[privKeysDictionary[j]] != j {
			fmt.Println("Test error!")
			break
		}
	}
	///end test---------------------------------------------------------------------

	var wordValue [wordsCount]int64
	for k := range wordValue {
		wordValue[k] = wordsDictionary[words[k]]
	}
	code := decode2Code(wordValue)

	if verifyChecksum(code) {
		key := new(big.Int).SetBytes(code[addressChecksum8bitsLen:])
		return key
	} else {
		return big.NewInt(0)
	}
}