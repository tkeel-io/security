package ciperutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func GenerateRandCipher(length int) string {
	return GenerateRandString(length)
}

func GenerateCipherKey32() string {
	return GenerateRandCipher(32)
}

func AESEncrypt(orig string, key string) (string, error) {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])

	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	// 使用 RawURLEncoding 不要使用 StdEncoding
	// 不要使用 StdEncoding  放在 url 参数中回导致错误
	// return base64.RawURLEncoding.EncodeToString(cryted)
	return base64.StdEncoding.EncodeToString(cryted), nil

}

func AESDecrypt(encrypted string, key string) (string, error) {
	// 使用 RawURLEncoding 不要使用 StdEncoding
	// 不要使用 StdEncoding  放在 url 参数中会导致错误
	// encryptedBytes, err := base64.RawURLEncoding.DecodeString(encrypted)
	encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		panic("cant decode data")
	}
	k := []byte(key)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])

	orig := make([]byte, len(encryptedBytes))
	// 解密
	blockMode.CryptBlocks(orig, encryptedBytes)
	// 去补全码
	orig = PKCS7UnPadding(orig)

	return string(orig), nil
}

// PKCS7Padding 补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
