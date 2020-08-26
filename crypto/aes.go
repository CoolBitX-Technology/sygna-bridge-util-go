package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

func aesEncrypt(plaintext, key, iv []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	blockSize := block.BlockSize()
	plaintext = pkcs7Padding(plaintext, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(plaintext))
	blockMode.CryptBlocks(crypted, plaintext)
	return crypted
}

func aesDecrypt(ciphertext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(origData, ciphertext)
	origData = pkcs7UnPadding(origData)
	return origData, nil
}
