package ciperutil

import (
	"crypto/sha1"
	"fmt"
)

func Hash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func HashCountWithBytes(data []byte) string {
	h := sha1.New()
	h.Write(data)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func CheckHash(expect string, data string) bool {
	return expect == Hash(data)
}

func CheckHashBytes(expect string, data []byte) bool {
	return expect == HashCountWithBytes(data)
}
