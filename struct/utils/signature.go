package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"log"
	"math/big"
)

// Signature represents the digital signature with R and S values.
type signature struct {
	R, S *big.Int
}

// Sign generates a digital signature for the given data using the private key.
func Sign(privateKey *ecdsa.PrivateKey, data interface{}) *Signature {
	// Marshal the data to JSON
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal data for signing: %v", err)
	}

	// Hash the data using SHA-256
	hash := sha256.Sum256(dataBytes)

	// Sign the hash using the private key
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		log.Fatalf("Failed to sign data: %v", err)
	}

	return &Signature{R: r, S: s}
}
