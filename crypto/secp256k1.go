package crypto

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/iancoleman/orderedmap"
)

//Sign Sign data with provided Private Key.
func Sign(message *orderedmap.OrderedMap, privateKey string) error {

	bPrivateKey, err := hex.DecodeString(privateKey)
	if err != nil {
		return err
	}

	message.Set("signature", "")

	bMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	bSignature, err := secp256k1.Sign(sha256Sum(bMessage), bPrivateKey)
	if err != nil {
		return err
	}

	signature := hex.EncodeToString(bSignature[:64])
	message.Set("signature", signature)

	return nil
}

//Verify Verify data with provided Public Key
func Verify(message *orderedmap.OrderedMap, publicKey string) (bool, error) {
	bPublicKey, err := hex.DecodeString(publicKey)
	if err != nil {
		return false, err
	}

	clone, err := cloneOrderedMap(message)
	if err != nil {
		return false, err
	}
	signature, exist := clone.Get("signature")
	if !exist {
		return false, errors.New("message must contain signature")
	}

	bSignature, err := hex.DecodeString(signature.(string))
	if err != nil {
		return false, err
	}

	clone.Set("signature", "")
	bClone, err := json.Marshal(clone)
	if err != nil {
		return false, err
	}

	return secp256k1.VerifySignature(bPublicKey, sha256Sum(bClone), bSignature), nil
}
