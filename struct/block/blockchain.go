package block

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	// Iterate through all blocks and transactions to collect wallet addresses
	for _, block := range bc.chain {
		for _, tx := range block.transactions {
			walletMap[tx.senderBlockchainAddress] = true
			walletMap[tx.recipientBlockchainAddress] = true
		}
	}

	// Convert map to slice
	wallets := make([]string, 0, len(walletMap))
	for wallet := range walletMap {
		wallets = append(wallets, wallet)
	}

	return wallets
}

// GetNeighbors returns the list of connected nodes
func (bc *Blockchain) GetNeighbors() []string {
	bc.muxNeighbors.Lock()
	defer bc.muxNeighbors.Unlock()
	return bc.neighbors
}
