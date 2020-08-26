package crypto

import (
	"bytes"
	"encoding/hex"
	"errors"

	"github.com/iancoleman/orderedmap"
)

// Encrypt Encrypt private info to hex string.
func Encrypt(sensitiveData []byte, publicKey string) (string, error) {
	eciesPublicKey, err := newPublicKeyFromHex(publicKey)
	if err != nil {
		return "", err
	}

	// Generate ephemeral key
	ek, err := generateKey()
	if err != nil {
		return "", err
	}

	encryptionKey, macKey, err := ek.Encapsulate(eciesPublicKey)

	if err != nil {
		return "", err
	}
	iv := make([]byte, 16)

	ciphertext := aesEncrypt(sensitiveData, encryptionKey, iv)
	dataToMac := appendBytes(iv, ek.publicKey.Bytes(false), ciphertext)

	encryptedData := appendBytes(ek.publicKey.Bytes(false), sha1Sum(dataToMac, macKey), ciphertext)
	encryptedHex := hex.EncodeToString(encryptedData)
	return encryptedHex, nil
}

//Decrypt Decrypt private info from recipient server.
func Decrypt(encryptedData, privateKey string) (interface{}, error) {
	bEncrypted, err := hex.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	eciesPrivateKey, err := newPrivateKeyFromHex(privateKey)
	if err != nil {
		return nil, err
	}
	ephemeralPubKey := bEncrypted[:65]
	mac := bEncrypted[65:85]
	ciphertext := bEncrypted[85:]

	eciesPublicKey, err := newPublicKeyFromBytes(ephemeralPubKey)
	if err != nil {
		return nil, err
	}

	encryptionKey, macKey, err := eciesPrivateKey.Encapsulate(eciesPublicKey)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, 16)

	dataToMac := appendBytes(iv, ephemeralPubKey, ciphertext)

	realMac := sha1Sum(dataToMac, macKey)

	if !bytes.Equal(realMac, mac) {
		return nil, errors.New("mac is not same")
	}

	decrypted, err := aesDecrypt(ciphertext, encryptionKey, iv)
	if err != nil {
		return nil, err
	}

	o := orderedmap.New()
	o.UnmarshalJSON(decrypted)

	if isOrderedMapEmpty(o) {
		return string(decrypted), nil
	}
	return o, nil
}
