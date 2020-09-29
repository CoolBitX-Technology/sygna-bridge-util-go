package crypto

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/fomichev/secp256k1"
)

// PrivateKey is an instance of secp256k1 private key with nested public key
type privateKey struct {
	*publicKey
	D *big.Int
}

// generateKey generates secp256k1 key pair
func generateKey() (*privateKey, error) {
	curve := secp256k1.SECP256K1()

	p, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("cannot generate key pair: %w", err)
	}

	return &privateKey{
		publicKey: &publicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(p),
	}, nil
}

// NewPrivateKeyFromHex decodes hex form of private key raw bytes, computes public key and returns PrivateKey instance
func newPrivateKeyFromHex(s string) (*privateKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("cannot decode hex string: %w", err)
	}

	return newPrivateKeyFromBytes(b), nil
}

// NewPrivateKeyFromBytes decodes private key raw bytes, computes public key and returns PrivateKey instance
func newPrivateKeyFromBytes(priv []byte) *privateKey {
	curve := secp256k1.SECP256K1()
	x, y := curve.ScalarBaseMult(priv)

	return &privateKey{
		publicKey: &publicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(priv),
	}
}

// Encapsulate encapsulates key by using Key Encapsulation Mechanism and returns symmetric key;
// can be safely used as encryption key
func (k *privateKey) Encapsulate(pub *publicKey) ([]byte, []byte, error) {
	if pub == nil {
		return nil, nil, fmt.Errorf("public key is empty")
	}
	sx, _ := pub.Curve.ScalarMult(pub.X, pub.Y, k.D.Bytes())

	hash := sha512Sum(sx.Bytes())
	return hash[:32], hash[32:], nil
}
