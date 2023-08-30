package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenKey(length int, data ...string) string {
	// generate only one string to apply sha256
	combinedData := ""
	for _, d := range data {
		combinedData += d
	}

	// hashing the string
	hash := sha256.Sum256([]byte(combinedData))
	hashString := hex.EncodeToString(hash[:])

	// Truncate the hash to desired length
	// ! This is only for the ease of use during the test,
	// ! then look for a better method
	if length > len(hashString) {
		length = len(hashString)
	}

	truncatedHash := hashString[:length]

	return truncatedHash
}
