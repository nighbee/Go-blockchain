package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *BlockchainServerHandler) Reset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.server.GetBlockchain().Reset()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Blockchain has been reset successfully",
	})
}
