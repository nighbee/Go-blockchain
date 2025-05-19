package handlers

import (
	"encoding/json"
	"net/http"
)

// Reset handles the reset request
func (h *BlockchainServerHandler) Reset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Reset the blockchain
	h.server.GetBlockchain().Reset()

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Blockchain has been reset successfully",
	})
}
