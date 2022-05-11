package ciperutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptAndDecrypt(t *testing.T) {
	key := GenerateCipherKey32()
	content := "secret"

	assert.Equal(t, 32, len([]byte(key)))

	e, err := AESEncrypt(content, key)
	assert.Nil(t, err)
	d, err := AESDecrypt(e, key)
	assert.Nil(t, err)
	assert.Equal(t, content, d)
}
