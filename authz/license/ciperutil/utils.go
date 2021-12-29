package ciperutil

import (
	"math/rand"
	"time"
)

var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func GenerateRandString(len int) string {
	b := make([]byte, len, len)
	for i := 0; i < len; i++ {
		rand.NewSource(time.Now().UnixNano())
		index := rand.Intn(len)
		b[i] = chars[index]
	}

	return string(b)
}
