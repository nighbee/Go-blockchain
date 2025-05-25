package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *BlockchainServerHandler) GetNodes(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch req.Method {
	case http.MethodGet:
		bc := h.server.GetBlockchain()
		neighbors := bc.GetNeighbors()

		response := struct {
			Nodes []string `json:"nodes"`
		}{
			Nodes: neighbors,
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
