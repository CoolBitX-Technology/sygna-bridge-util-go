package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"github.com/iancoleman/orderedmap"
)

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func sha1Sum(data, key []byte) []byte {
	hash := hmac.New(sha1.New, key)
	hash.Write(data)

	return hash.Sum(nil)
}

func sha256Sum(msg []byte) []byte {
	hash := sha256.New()
	hash.Write(msg)

	return hash.Sum(nil)
}

func sha512Sum(secret []byte) [sha512.Size]byte {
	hash := sha512.Sum512(secret)

	return hash
}

func appendBytes(data ...[]byte) []byte {
	result := make([]byte, 0)
	for _, v := range data {
		result = append(result, v...)
	}
	return result
}

func isOrderedMapEmpty(o *orderedmap.OrderedMap) bool {
	return len(o.Keys()) == 0
}

func cloneOrderedMap(o *orderedmap.OrderedMap) (*orderedmap.OrderedMap, error) {
	b, err := o.MarshalJSON()
	if err != nil {
		return nil, err
	}
	clone := orderedmap.New()
	err = clone.UnmarshalJSON(b)
	if err != nil {
		return nil, err
	}
	return clone, nil
}
