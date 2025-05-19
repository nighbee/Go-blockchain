package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *BlockchainServerHandler) GetWallets(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		// Get all wallets from the blockchain
		wallets := h.server.GetBlockchain().GetWallets()

		// Prepare response
		response := struct {
			Wallets []string `json:"wallets"`
		}{
			Wallets: wallets,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("ERROR: Failed to encode response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		log.Printf("ERROR: Invalid HTTP Method: %s", req.Method)
		w.WriteHeader(http.StatusBadRequest)
	}
}
