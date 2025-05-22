package handlers

import (
	"block/struct/utils"
	"io"
	"log"
	"net/http"
)

func (h *BlockchainServerHandler) Consensus(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:

		bc := h.server.GetBlockchain()

		replaced := bc.ResolveConflicts()

		w.Header().Add("Content-Type", "application/json")

		if replaced {
			io.WriteString(w, string(utils.JsonStatus("success")))
		} else {
			io.WriteString(w, string(utils.JsonStatus("fail")))
		}
	default:

		log.Printf("ERROR: Invalid HTTP Method")

		w.WriteHeader(http.StatusBadRequest)
	}
}
