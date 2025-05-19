package handlers

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (h *BlockchainServerHandler) RegisterWallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		// Read the request body
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("Failed to read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Restore the request body to its original state
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// Define a struct to capture the request body
		type RequestBody struct {
			BlockchainAddress *string `json:"blockchainAddress"`
		}

		// Decode the request body into the struct
		var requestBody RequestBody
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			log.Printf("ERROR: Failed to decode request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Generate ECDSA key pair
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			log.Printf("ERROR: Failed to generate private key: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		publicKey := privateKey.PublicKey
		// Ensure public key is 128 characters (64 bytes: 32 for X, 32 for Y)
		publicKeyBytes := elliptic.Marshal(elliptic.P256(), publicKey.X, publicKey.Y)
		publicKeyHex := hex.EncodeToString(publicKeyBytes[1:]) // Skip leading 04
		if len(publicKeyHex) != 128 {
			log.Printf("ERROR: Generated public key length %d, expected 128", len(publicKeyHex))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		privateKeyHex := fmt.Sprintf("%064x", privateKey.D.Bytes())

		// Use provided blockchainAddress or generate one
		var address string
		if requestBody.BlockchainAddress != nil && *requestBody.BlockchainAddress != "" {
			address = *requestBody.BlockchainAddress
		} else {
			// Generate address from public key (SHA256 hash prefix)
			hash := sha256.Sum256([]byte(publicKeyHex))
			address = "1" + hex.EncodeToString(hash[:20])
		}

		// Register the wallet
		success := h.server.GetBlockchain().RegisterNewWallet(address, "REGISTER USER WALLET")
		if !success {
			log.Printf("ERROR: Failed to register wallet with address: %s", address)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Prepare response
		response := struct {
			Address    string `json:"address"`
			PublicKey  string `json:"public_key"`
			PrivateKey string `json:"private_key"`
		}{
			Address:    address,
			PublicKey:  publicKeyHex,
			PrivateKey: privateKeyHex,
		}

		log.Printf("private_key %s", privateKeyHex)
		log.Printf("public_key %s", publicKeyHex)
		log.Printf("address %s", address)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("ERROR: Failed to encode response: %v", err)
		}

	default:
		log.Printf("ERROR: Invalid HTTP Method: %s", req.Method)
		w.WriteHeader(http.StatusBadRequest)
	}
}
