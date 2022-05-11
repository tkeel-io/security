package ciperutil

import (
	"crypto/rand"
	"math/big"
)

var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func GenerateRandString(length int) string {
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		bigIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(length)))
		index := bigIndex.Int64()
		b[i] = chars[index]
	}

	return string(b)
}
