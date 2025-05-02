package handlers

import (
	"block/struct/block"
	"block/struct/wallet"
	"fmt"
	"reflect"
)

func LogMethods(i interface{}) {
	t := reflect.TypeOf(i)

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fmt.Println(method.Name)
	}
}

type BlockchainServer interface {
	Port() uint16
	GetWallet() *wallet.Wallet
	GetBlockchain() *block.Blockchain
	// TODO: Learn interfaces
	// TODO: Add GetBlockchain() method into controllers
}

type BlockchainServerHandler struct {
	server BlockchainServer
}

func NewBlockchainServerHandler(s BlockchainServer) *BlockchainServerHandler {
	// LogMethods(s)
	return &BlockchainServerHandler{server: s}
}
