package handlers

import (
	"block/struct/block"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func (h *BlockchainServerHandler) Balance(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:

		br := &block.BalanceResponse{} // Use the BalanceResponse type

		blockchainAddress := req.URL.Query().Get("blockchainAddress")

		balance, err := h.server.GetBlockchain().CalculateTotalBalance(blockchainAddress)

		br.Balance = balance
		if err != nil {
			log.Printf("ERROR: %v", err)
			br.Error = err.Error()
		}

		m, _ := json.Marshal(br)

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(m))

	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}
