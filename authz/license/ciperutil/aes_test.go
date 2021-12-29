package ciperutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
	key := GenerateCipherKey32()
	content := "secret"

	assert.Equal(t, 32, len([]byte(key)))

	e := AesEncrypt(content, key)
	d := AesDecrypt(e, key)
	assert.Equal(t, content, d)
}
