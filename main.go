package bridgeutil

import (
	"encoding/json"

	"github.com/CoolBitX-Technology/sygna-bridge-util-go/crypto"
	"github.com/iancoleman/orderedmap"
)

//Encrypt Encrypt private info to hex string.
func Encrypt(sensitiveData *orderedmap.OrderedMap, publicKey string) (string, error) {
	b, err := json.Marshal(sensitiveData)
	if err != nil {
		return "", err
	}
	return crypto.Encrypt(b, publicKey)
}

//EncryptString Encrypt private info(string) to hex string.
func EncryptString(sensitiveData, publicKey string) (string, error) {
	b := []byte(sensitiveData)
	return crypto.Encrypt(b, publicKey)
}

//Decrypt Decrypt private info from recipient server.
func Decrypt(encryptedData, privateKey string) (interface{}, error) {
	return crypto.Decrypt(encryptedData, privateKey)
}

//Sign Sign data with provided Private Key.
func Sign(message *orderedmap.OrderedMap, privateKey string) error {
	return crypto.Sign(message, privateKey)
}

//Verify Verify data with provided Public Key or default sygna bridge
func Verify(message *orderedmap.OrderedMap, publicKey ...string) (bool, error) {
	defaultPublicKey := SygnaBridgeCentralPubkey
	if len(publicKey) > 0 {
		defaultPublicKey = publicKey[0]
	}
	return crypto.Verify(message, defaultPublicKey)
}
