package block

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// ==============================
// Blockchain Struct and Methods
// ==============================

// Blockchain represents the entire blockchain structure.
type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
	port              uint16
	mux               sync.Mutex
	neighbors         []string
	muxNeighbors      sync.Mutex
}

// NewBlockchain creates a new instance of Blockchain.
func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock([]*Transaction{}, b.GetHash())
	bc.port = port
	return bc
}

// Chain returns the chain of the Blockchain.
func (bc *Blockchain) Chain() []*Block {
	return bc.chain
}

// Run initializes and runs the Blockchain.
func (bc *Blockchain) Run() {
	bc.StartSyncNeighbors()
	bc.ResolveConflicts()
	bc.StartMining() // Start mining automatically
}

func (bc *Blockchain) CreateBlock(transactions []*Transaction, previousHash string) *Block {
	b := NewBlock(transactions, previousHash)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	for _, n := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/transactions", n)
		client := &http.Client{}
		req, _ := http.NewRequest("DELETE", endpoint, nil)
		resp, _ := client.Do(req)
		log.Printf("%v", resp)
	}
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) GetBlocks(amount int) []*Block {
	n := len(bc.chain)
	var last10Blocks []*Block
	if n > amount {
		last10Blocks = append([]*Block(nil), bc.chain[n-amount:n]...)
	} else {
		last10Blocks = append([]*Block(nil), bc.chain...)
	}

	// Reverse the slice
	for i := len(last10Blocks)/2 - 1; i >= 0; i-- {
		opp := len(last10Blocks) - 1 - i
		last10Blocks[i], last10Blocks[opp] = last10Blocks[opp], last10Blocks[i]
	}

	return last10Blocks
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i,
			strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chain"`
	}{
		Blocks: bc.chain,
	})
}

func (bc *Blockchain) UnmarshalJSON(data []byte) error {
	v := &struct {
		Blocks *[]*Block `json:"chain"`
	}{
		Blocks: &bc.chain,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

// MineBlock creates a new block with pending transactions and rewards the miner
func (bc *Blockchain) MineBlock(minerAddress string) (bool, error) {
	// Lock the blockchain while mining
	bc.mux.Lock()
	defer bc.mux.Unlock()

	// Get the last block
	lastBlock := bc.chain[len(bc.chain)-1]

	// Create a new block with current transactions
	newBlock := NewBlock(bc.transactionPool, lastBlock.GetHash())

	// Add mining reward transaction
	rewardTx := &Transaction{
		senderBlockchainAddress:    MINING_SENDER,
		recipientBlockchainAddress: minerAddress,
		value:                      MINING_REWARD,
		message:                    "MINING REWARD",
	}
	newBlock.transactions = append(newBlock.transactions, rewardTx)

	// Mine the block (find a valid hash)
	for {
		hash := newBlock.CalculateHash()
		newBlock.SetHash(hash)
		if newBlock.IsValidHash() {
			break
		}
		newBlock.SetNonce(newBlock.GetNonce() + 1)
	}

	// Add the new block to the chain
	bc.chain = append(bc.chain, newBlock)

	// Clear the transaction pool
	bc.transactionPool = []*Transaction{}

	return true, nil
}

// GetWallets returns a list of all registered wallet addresses in the blockchain
func (bc *Blockchain) GetWallets() []string {
	// Use a map to store unique wallet addresses
	walletMap := make(map[string]bool)

	log.Printf("Starting to scan blockchain for registered wallets...")
	log.Printf("Total blocks in chain: %d", len(bc.chain))

	// Iterate through all blocks and transactions to collect wallet addresses
	for blockIndex, block := range bc.chain {
		log.Printf("Scanning block #%d with %d transactions", blockIndex, len(block.transactions))

		for txIndex, tx := range block.transactions {
			// Debug logging for each transaction
			log.Printf("Transaction #%d in block #%d:", txIndex, blockIndex)
			log.Printf("  From: %s", tx.senderBlockchainAddress)
			log.Printf("  To: %s", tx.recipientBlockchainAddress)
			log.Printf("  Message: %s", tx.message)
			log.Printf("  Amount: %f", tx.value)

			// Only include wallets that were registered with "REGISTER USER WALLET" message
			if tx.message == "REGISTER USER WALLET" &&
				tx.senderBlockchainAddress == "THE BLOCKCHAIN" {
				// Check if wallet is already in the map
				if !walletMap[tx.recipientBlockchainAddress] {
					walletMap[tx.recipientBlockchainAddress] = true
					log.Printf("  Found new registered wallet: %s", tx.recipientBlockchainAddress)
				} else {
					log.Printf("  Skipping duplicate wallet: %s", tx.recipientBlockchainAddress)
				}
			}
		}
	}

	// Convert map to slice and sort for consistent display
	wallets := make([]string, 0, len(walletMap))
	for wallet := range walletMap {
		wallets = append(wallets, wallet)
	}

	// Sort wallets for consistent display
	sort.Strings(wallets)

	log.Printf("Total unique registered wallets found: %d", len(wallets))
	for i, wallet := range wallets {
		log.Printf("Registered wallet %d: %s", i+1, wallet)
	}

	return wallets
}

// GetNeighbors returns the list of connected nodes
func (bc *Blockchain) GetNeighbors() []string {
	bc.muxNeighbors.Lock()
	defer bc.muxNeighbors.Unlock()
	return bc.neighbors
}

// Reset resets the blockchain to its initial state
func (bc *Blockchain) Reset() {
	// Lock the blockchain while resetting
	bc.mux.Lock()
	defer bc.mux.Unlock()

	log.Printf("Starting blockchain reset...")

	// Clear the chain and transaction pool
	bc.chain = []*Block{}
	bc.transactionPool = []*Transaction{}

	// Create a new genesis block
	b := &Block{}
	bc.CreateBlock([]*Transaction{}, b.GetHash())

	// Clear neighbors
	bc.muxNeighbors.Lock()
	bc.neighbors = []string{}
	bc.muxNeighbors.Unlock()

	// Delete the blockchain file
	filePath := filepath.Join("data", blockchainFile)
	if err := os.Remove(filePath); err != nil {
		log.Printf("Warning: Could not delete blockchain file: %v", err)
	}

	// Save the reset blockchain
	if err := bc.SaveBlockchain(); err != nil {
		log.Printf("ERROR: Failed to save reset blockchain: %v", err)
	} else {
		log.Printf("Blockchain successfully reset and saved")
	}
}

func (bc *Blockchain) ResolveConflicts() bool {
	// Initialize variables to track the longest chain and its length
	var longestChain []*Block = nil
	maxLength := len(bc.chain)

	for _, n := range bc.neighbors {
		fmt.Println("Resolve conflict with ", n)

		endpoint := fmt.Sprintf("%s/chain", n)

		resp, err := http.Get(endpoint)
		if err != nil {
			log.Printf("ERROR: Failed to fetch chain from neighbor %s: %v", n, err)
			continue // Skip to the next neighbor in case of error
		}

		if resp.StatusCode == http.StatusOK {
			var bcResp Blockchain
			decoder := json.NewDecoder(resp.Body)

			err := decoder.Decode(&bcResp)
			if err != nil {
				log.Printf("ERROR: Failed to decode JSON response from neighbor %s: %v", n, err)
				continue // Skip to the next neighbor in case of error
			}

			chain := bcResp.Chain()

			if len(chain) > maxLength && bc.ValidChain(chain) {
				maxLength = len(chain)
				longestChain = chain
			}
		} else {
			log.Printf("WARNING: Failed to fetch chain from neighbor %s. Status code: %d", n, resp.StatusCode)
		}
	}

	if longestChain != nil {
		bc.chain = longestChain
		log.Printf("INFO: Resolved conflicts. Replaced blockchain with the longest valid chain.")

		// Save blockchain after resolving conflicts
		if err := bc.SaveBlockchain(); err != nil {
			log.Printf("ERROR: Failed to save blockchain after resolving conflicts: %v", err)
		}

		return true
	}

	log.Printf("INFO: No longer valid chain found among neighbors. No conflicts resolved.")
	return false
}
