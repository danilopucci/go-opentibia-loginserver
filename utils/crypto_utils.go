package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func Sha1Hash(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes) // Convert to a hex string
}
