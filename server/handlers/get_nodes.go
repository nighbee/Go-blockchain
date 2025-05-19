package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *BlockchainServerHandler) GetNodes(w http.ResponseWriter, req *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight request
	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch req.Method {
	case http.MethodGet:
		// Get the blockchain instance
		bc := h.server.GetBlockchain()

		// Get the list of neighbors
		neighbors := bc.GetNeighbors()

		// Prepare response
		response := struct {
			Nodes []string `json:"nodes"`
		}{
			Nodes: neighbors,
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")

		// Encode and send response
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
