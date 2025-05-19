package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// ==============================
// Constants
// ==============================

// Mining related constants.
const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
	MINING_TIMER_SEC  = 20
)

// Network related constants.
const (
	BLOCKCHAIN_PORT_RANGE_START      = 5001
	BLOCKCHAIN_PORT_RANGE_END        = 5002
	NEIGHBOR_IP_RANGE_START          = 0
	NEIGHBOR_IP_RANGE_END            = 1
	BLOCKCHIN_NEIGHBOR_SYNC_TIME_SEC = 20
)

// Block represents a block in the blockchain.
type Block struct {
	timestamp    int64
	transactions []*Transaction
	prevHash     string
	hash         string
	nonce        int
}

// NewBlock creates a new block with the given parameters.
func NewBlock(transactions []*Transaction, prevHash string) *Block {
	return &Block{
		timestamp:    time.Now().Unix(),
		transactions: transactions,
		prevHash:     prevHash,
		hash:         "",
		nonce:        0,
	}
}

// Accessor methods for the Block attributes.
func (b *Block) CalculateHash() string {
	m, err := json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
		PrevHash     string         `json:"prevHash"`
		Nonce        int            `json:"nonce"`
	}{
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
		PrevHash:     b.prevHash,
		Nonce:        b.nonce,
	})
	if err != nil {
		log.Printf("ERROR: Failed to marshal block: %v", err)
		return ""
	}

	hash := sha256.Sum256(m)
	return fmt.Sprintf("%x", hash)
}

func (b *Block) IsValidHash() bool {
	return b.hash[:2] == "00"
}

func (b *Block) GetHash() string {
	return b.hash
}

func (b *Block) GetPrevHash() string {
	return b.prevHash
}

func (b *Block) GetTransactions() []*Transaction {
	return b.transactions
}

func (b *Block) GetTimestamp() int64 {
	return b.timestamp
}

func (b *Block) GetNonce() int {
	return b.nonce
}

func (b *Block) SetHash(hash string) {
	b.hash = hash
}

func (b *Block) SetNonce(nonce int) {
	b.nonce = nonce
}

// Print displays the block's attributes.
func (b *Block) Print() {
	fmt.Printf("timestamp       %d\n", b.timestamp)
	fmt.Printf("nonce           %d\n", b.nonce)
	fmt.Printf("previousHash   %s\n", b.prevHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

// JSON Handling for Block

// MarshalJSON customizes the JSON encoding of the block.
func (b *Block) MarshalJSON() ([]byte, error) {
	transactions := b.transactions
	if transactions == nil {
		transactions = []*Transaction{}
	}
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
		PrevHash     string         `json:"prevHash"`
		Nonce        int            `json:"nonce"`
	}{
		Timestamp:    b.timestamp,
		Transactions: transactions,
		PrevHash:     b.prevHash,
		Nonce:        b.nonce,
	})
}

// UnmarshalJSON customizes the JSON decoding of the block.
func (b *Block) UnmarshalJSON(data []byte) error {
	v := &struct {
		Timestamp    *int64          `json:"timestamp"`
		Transactions *[]*Transaction `json:"transactions"`
		PrevHash     *string         `json:"prevHash"`
		Nonce        *int            `json:"nonce"`
	}{
		Timestamp:    &b.timestamp,
		Transactions: &b.transactions,
		PrevHash:     &b.prevHash,
		Nonce:        &b.nonce,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}
