package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func (h *BlockchainServerHandler) GetBlocks(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := h.server.GetBlockchain()
		m, _ := json.Marshal(bc.GetBlocks(10))
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}
