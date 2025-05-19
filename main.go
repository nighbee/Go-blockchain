package main

import (
	"block/middleware"
	"block/server/handlers"
	"block/struct/block"
	"block/struct/utils"
	"block/struct/wallet"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

var cache map[string]*block.Blockchain = make(map[string]*block.Blockchain)

type BlockchainServer struct {
	port   uint16
	Wallet *wallet.Wallet
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) GetWallet() *wallet.Wallet {
	return bcs.Wallet
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{
		port:   port,
		Wallet: nil,
	}
}

func (bcs *BlockchainServer) GetBlockchain() *block.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		// Try to load existing blockchain first
		bc = block.LoadBlockchain(bcs.Port())
		if bc == nil {
			// If no existing blockchain found, create a new one
			minersWallet := wallet.NewWallet()
			bc = block.NewBlockchain(minersWallet.BlockchainAddress(), bcs.Port())
			bcs.Wallet = minersWallet

			log.Printf("Created new blockchain")
			log.Printf("privateKey %v", minersWallet.PrivateKeyStr())
			log.Printf("publicKey %v", minersWallet.PublicKeyStr())
			log.Printf("blockchainAddress %v", minersWallet.BlockchainAddress())

			// Create a sender wallet for transactions
			senderWallet := wallet.NewWallet()

			// Add real transactions and mine blocks
			transaction1 := block.NewTransaction("Payment for services", "recipient_1", "sender_1", 25)
			signature1 := utils.Sign(senderWallet.PrivateKey(), transaction1)
			bc.AddTransaction("sender_1", "recipient_1", "Payment for services", 25, senderWallet.PublicKey(), signature1)
			bc.Mining()

			transaction2 := block.NewTransaction("Payment for goods", "recipient_2", "sender_2", 20.0)
			signature2 := utils.Sign(senderWallet.PrivateKey(), transaction2)
			bc.AddTransaction("sender_2", "recipient_2", "Payment for goods", 20.0, senderWallet.PublicKey(), signature2)
			bc.Mining()

			transaction3 := block.NewTransaction("Payment for rent", "recipient_3", "sender_3", 30)
			signature3 := utils.Sign(senderWallet.PrivateKey(), transaction3)
			bc.AddTransaction("sender_3", "recipient_3", "Payment for rent", 30, senderWallet.PublicKey(), signature3)
			bc.Mining()

			// Save the new blockchain
			if err := bc.SaveBlockchain(); err != nil {
				log.Printf("Error saving new blockchain: %v", err)
			}
		} else {
			log.Printf("Loaded existing blockchain")
		}
		cache["blockchain"] = bc
	}
	return bc
}

// enableCORS is a middleware function that adds CORS headers to the response.
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // Allowed HTTP methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")     // Allowed headers
		w.Header().Set("Access-Control-Allow-Credentials", "true")                        // Allow credentials

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (bcs *BlockchainServer) Run() {
	bcs.GetBlockchain().Run()

	router := mux.NewRouter()
	router.Use(utils.CorsMiddleware())

	handler := handlers.NewBlockchainServerHandler(bcs)

	router.HandleFunc("/chain", handler.GetChain)
	router.HandleFunc("/balance", handler.Balance)
	router.HandleFunc("/consensus", handler.Consensus)
	router.HandleFunc("/mine", handler.HandleMine)
	router.HandleFunc("/mine/start", handler.StartMine)
	router.HandleFunc("/miner/blocks", handler.GetBlocks)
	router.HandleFunc("/miner/wallet", handler.MinerWallet)
	router.HandleFunc("/transactions", handler.Transactions)
	router.HandleFunc("/wallet/register", handler.RegisterWallet)
	router.HandleFunc("/wallets", handler.GetWallets)
	router.HandleFunc("/nodes", handler.GetNodes)
	router.HandleFunc("/reset", handler.Reset)
	router.HandleFunc("/sign", handler.HandleSign).Methods("POST")
	// Chain middlewares: enableCORS -> LoggingMiddleware
	corsAndLoggingHandler := middleware.LoggingMiddleware(enableCORS(router))

	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.Port())), corsAndLoggingHandler))
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		port = 5001
	}

	log.Printf("Port: %d\n", port)

	app := NewBlockchainServer(uint16(port))
	app.Run()
}
