package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashAadhar(aadhar string) string {
	hasher := sha256.New()
	hasher.Write([]byte(aadhar))
	return hex.EncodeToString(hasher.Sum(nil))
}
