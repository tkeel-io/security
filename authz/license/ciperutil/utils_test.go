package ciperutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateRandString(t *testing.T) {
	s1 := GenerateRandString(10)
	s2 := GenerateRandString(10)
	s3 := GenerateRandString(10)
	s4 := GenerateRandString(10)
	assert.Equal(t, 10, len(s1))
	assert.Equal(t, 10, len(s2))
	assert.Equal(t, 10, len(s3))
	assert.Equal(t, 10, len(s4))
	assert.NotEqual(t, s1, s2)
	assert.NotEqual(t, s3, s4)

	s01 := GenerateCipherKey32()
	s02 := GenerateCipherKey32()
	assert.NotEqual(t, s01, s02)
}

func TestAesEncryptAndDecrypt(t *testing.T) {
	cipher := GenerateCipherKey32()
	data := "secret_data"
	en := AesEncrypt(data, cipher)
	assert.Equal(t, data, AesDecrypt(en, cipher))
}
