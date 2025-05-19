package handlers

import (
	"log"
	"net/http"
)

func (h *BlockchainServerHandler) Reset(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		// Reset the blockchain
		h.server.GetBlockchain().Reset()

		// Return success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Blockchain reset successfully"))

	default:
		log.Printf("ERROR: Invalid HTTP Method: %s", req.Method)
		w.WriteHeader(http.StatusBadRequest)
	}
}
