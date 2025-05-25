package handlers

import (
	"block/struct/block"
	"block/struct/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type TransactionRequest struct {
	SenderBlockchainAddress    string  `json:"senderBlockchainAddress"`
	RecipientBlockchainAddress string  `json:"recipientBlockchainAddress"`
	Message                    string  `json:"message"`
	Value                      float32 `json:"value"`
	SenderPublicKey            string  `json:"senderPublicKey"`
	Signature                  string  `json:"signature"`
}

func (h *BlockchainServerHandler) HandleGetTransaction(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bc := h.server.GetBlockchain()
	transactions := bc.TransactionPool()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transactions": transactions,
		"length":       len(transactions),
	})
}

func (h *BlockchainServerHandler) HandlePostTransaction(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var txReq TransactionRequest
	if err := json.NewDecoder(req.Body).Decode(&txReq); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Failed to decode request body: %v", err)})
		return
	}

	if txReq.SenderBlockchainAddress == "" || txReq.RecipientBlockchainAddress == "" ||
		txReq.Message == "" || txReq.Value <= 0 || txReq.SenderPublicKey == "" || txReq.Signature == "" {
		log.Printf("ERROR: Missing required fields")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Missing required fields"})
		return
	}

	publicKey, err := utils.PublicKeyFromString(txReq.SenderPublicKey)
	if err != nil {
		log.Printf("ERROR: Invalid public key: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Invalid public key: %v", err)})
		return
	}

	signature, err := utils.SignatureFromString(txReq.Signature)
	if err != nil {
		log.Printf("ERROR: Invalid signature: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Invalid signature: %v", err)})
		return
	}

	success, err := h.server.GetBlockchain().AddTransaction(
		txReq.SenderBlockchainAddress,
		txReq.RecipientBlockchainAddress,
		txReq.Message,
		txReq.Value,
		publicKey,
		signature,
	)

	if !success {
		log.Printf("ERROR: Failed to create transaction: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Failed to create transaction: %v", err)})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Transaction created successfully",
		"status":  "success",
	})
}

func (h *BlockchainServerHandler) HandlePutTransaction(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t block.TransactionRequest
	err := decoder.Decode(&t)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, string(utils.JsonStatus("fail")))
		return
	}

	if !t.Validate() {
		log.Println("ERROR: missing field(s)")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, string(utils.JsonStatus("fail")))
		return
	}

	log.Printf("Received senderPublicKey: %s", *t.SenderPublicKey)
	log.Printf("Received signature: %s", *t.Signature)

	publicKey, err := utils.PublicKeyFromString(*t.SenderPublicKey)
	if err != nil {
		log.Printf("ERROR: Invalid sender public key: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "fail", "error": fmt.Sprintf("Invalid sender public key: %v", err)})
		return
	}

	signature, err := utils.SignatureFromString(*t.Signature)
	if err != nil {
		log.Printf("ERROR: Invalid signature: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "fail", "error": fmt.Sprintf("Invalid signature: %v", err)})
		return
	}

	bc := h.server.GetBlockchain()

	isUpdated, err := bc.AddTransaction(
		*t.SenderBlockchainAddress,
		*t.RecipientBlockchainAddress,
		*t.Message,
		*t.Value,
		publicKey,
		signature,
	)

	w.Header().Add("Content-Type", "application/json")

	var m []byte
	if err != nil {
		log.Printf("ERROR: Failed to add transaction: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "fail", "error": err.Error()})
		return
	}

	if !isUpdated {
		w.WriteHeader(http.StatusBadRequest)
		m = utils.JsonStatus("fail")
	} else {
		m = utils.JsonStatus("success")
	}

	io.WriteString(w, string(m))
}

func (h *BlockchainServerHandler) HandleDeleteTransaction(w http.ResponseWriter, req *http.Request) {
	bc := h.server.GetBlockchain()

	bc.ClearTransactionPool()

	io.WriteString(w, string(utils.JsonStatus("success")))
}

func (h *BlockchainServerHandler) Transactions(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		h.HandleGetTransaction(w, req)
	case http.MethodPost:
		h.HandlePostTransaction(w, req)
	case http.MethodPut:
		h.HandlePutTransaction(w, req)
	case http.MethodDelete:
		h.HandleDeleteTransaction(w, req)
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}
