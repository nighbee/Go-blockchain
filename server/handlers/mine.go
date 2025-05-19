package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

const MINING_REWARD = 1.0

type MineRequest struct {
	MinerAddress string `json:"minerAddress"`
}

func (h *BlockchainServerHandler) HandleMine(w http.ResponseWriter, req *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight request
	if req.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var mineReq MineRequest
	if err := json.NewDecoder(req.Body).Decode(&mineReq); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if mineReq.MinerAddress == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Miner address is required"})
		return
	}

	// Get the blockchain instance
	bc := h.server.GetBlockchain()

	// Mine a new block
	success, err := bc.MineBlock(mineReq.MinerAddress)
	if err != nil {
		log.Printf("ERROR: Mining failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if !success {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Mining failed"})
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Block mined successfully",
		"reward":  MINING_REWARD,
	})
}
