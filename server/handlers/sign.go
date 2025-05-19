package handlers

import (
	"block/struct/block"
	"block/struct/utils"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SignRequest struct {
	SenderBlockchainAddress    string  `json:"senderBlockchainAddress"`
	RecipientBlockchainAddress string  `json:"recipientBlockchainAddress"`
	Message                    string  `json:"message"`
	Value                      float32 `json:"value"`
	PrivateKey                 string  `json:"privateKey"`
	PublicKey                  string  `json:"publicKey"`
}

type SignResponse struct {
	Signature string `json:"signature"`
}

func (h *BlockchainServerHandler) HandleSign(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var signReq SignRequest
	if err := json.NewDecoder(req.Body).Decode(&signReq); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Failed to decode request body: %v", err)})
		return
	}

	// Validate required fields
	if signReq.SenderBlockchainAddress == "" {
		log.Printf("ERROR: Missing senderBlockchainAddress")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Missing senderBlockchainAddress"})
		return
	}
	if signReq.RecipientBlockchainAddress == "" {
		log.Printf("ERROR: Missing recipientBlockchainAddress")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Missing recipientBlockchainAddress"})
		return
	}
	if signReq.Message == "" {
		log.Printf("ERROR: Missing message")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Missing message"})
		return
	}
	if signReq.Value <= 0 {
		log.Printf("ERROR: Invalid value: %f", signReq.Value)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid value"})
		return
	}
	if signReq.PrivateKey == "" {
		log.Printf("ERROR: Missing privateKey")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Missing privateKey"})
		return
	}
	if signReq.PublicKey == "" {
		log.Printf("ERROR: Missing publicKey")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Missing publicKey"})
		return
	}

	// Log the request details
	log.Printf("Signing transaction: sender=%s, recipient=%s, message=%s, value=%f",
		signReq.SenderBlockchainAddress,
		signReq.RecipientBlockchainAddress,
		signReq.Message,
		signReq.Value)

	// Create transaction
	tx := block.NewTransaction(
		signReq.SenderBlockchainAddress,
		signReq.RecipientBlockchainAddress,
		signReq.Message,
		signReq.Value,
	)

	// Parse public key
	publicKey, err := utils.PublicKeyFromString(signReq.PublicKey)
	if err != nil {
		log.Printf("ERROR: Invalid public key: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Invalid public key: %v", err)})
		return
	}

	// Parse private key
	privateKey, err := utils.PrivateKeyFromString(signReq.PrivateKey, publicKey)
	if err != nil {
		log.Printf("ERROR: Invalid private key: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Invalid private key: %v", err)})
		return
	}

	// Sign the transaction
	m, err := json.Marshal(tx)
	if err != nil {
		log.Printf("ERROR: Failed to marshal transaction: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Failed to marshal transaction: %v", err)})
		return
	}
	log.Printf("Sign: Marshaled transaction: %s", string(m))

	// Create a hash of the transaction data
	hash := sha256.Sum256(m)
	log.Printf("Sign: Transaction hash: %x", hash)

	// Sign the hash with the private key
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		log.Printf("ERROR: Failed to sign transaction: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Failed to sign transaction: %v", err)})
		return
	}

	// Create signature object
	signature := &utils.Signature{R: r, S: s}

	// Log the signature details for debugging
	log.Printf("Generated signature R: %x", r)
	log.Printf("Generated signature S: %x", s)
	log.Printf("Generated signature string: %s", signature.String())

	// Verify the signature immediately after generation
	valid := ecdsa.Verify(&privateKey.PublicKey, hash[:], r, s)
	if !valid {
		log.Printf("ERROR: Generated signature failed immediate verification")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Generated signature failed verification"})
		return
	}
	log.Printf("Generated signature passed immediate verification")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SignResponse{Signature: signature.String()})
}
