package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *BlockchainServerHandler) GetWallets(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		wallets := h.server.GetBlockchain().GetWallets()

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
