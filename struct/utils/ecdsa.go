package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%064x%064x", s.R, s.S)
}

// String2BigIntTuple converts a hexadecimal string into a tuple of big integers.
// The input string must be a 128-character hex string (64 bytes when decoded).
func String2BigIntTuple(s string) (*big.Int, *big.Int, error) {
	if len(s) != 128 {
		return nil, nil, fmt.Errorf("invalid public key length: got %d characters, expected 128", len(s))
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode hex string: %v", err)
	}
	if len(b) != 64 {
		return nil, nil, fmt.Errorf("invalid decoded public key length: got %d bytes, expected 64", len(b))
	}
	x := new(big.Int).SetBytes(b[:32])
	y := new(big.Int).SetBytes(b[32:64])
	return x, y, nil
}

func SignatureFromString(s string) (*Signature, error) {
	if len(s) != 128 {
		return nil, fmt.Errorf("invalid signature length: got %d characters, expected 128", len(s))
	}
	x, y, err := String2BigIntTuple(s)
	if err != nil {
		return nil, err
	}
	return &Signature{R: x, S: y}, nil
}

// PublicKeyFromString converts a hex string to an ECDSA public key.
func PublicKeyFromString(s string) (*ecdsa.PublicKey, error) {
	x, y, err := String2BigIntTuple(s)
	if err != nil {
		return nil, err
	}
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, nil
}

// PrivateKeyFromString converts a hex string to an ECDSA private key.
func PrivateKeyFromString(s string, publicKey *ecdsa.PublicKey) (*ecdsa.PrivateKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}
	var bi big.Int
	bi.SetBytes(b)
	return &ecdsa.PrivateKey{PublicKey: *publicKey, D: &bi}, nil
}
