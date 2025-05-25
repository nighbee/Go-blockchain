package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"log"
	"math/big"
)

type signature struct {
	R, S *big.Int
}

func Sign(privateKey *ecdsa.PrivateKey, data interface{}) *Signature {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal data for signing: %v", err)
	}

	hash := sha256.Sum256(dataBytes)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		log.Fatalf("Failed to sign data: %v", err)
	}

	return &Signature{R: r, S: s}
}
