package block

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const blockchainFile = "blockchain.json"

// SaveBlockchain saves the blockchain to a file
func (bc *Blockchain) SaveBlockchain() error {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	// Marshal blockchain to JSON
	data, err := json.Marshal(bc)
	if err != nil {
		return fmt.Errorf("failed to marshal blockchain: %v", err)
	}

	// Write to file
	filePath := filepath.Join("data", blockchainFile)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write blockchain file: %v", err)
	}

	return nil
}

// LoadBlockchain loads the blockchain from a file
func LoadBlockchain(port uint16) *Blockchain {
	filePath := filepath.Join("data", blockchainFile)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading blockchain file: %v\n", err)
		return nil
	}

	// Unmarshal JSON
	var bc Blockchain
	if err := json.Unmarshal(data, &bc); err != nil {
		fmt.Printf("Error unmarshaling blockchain: %v\n", err)
		return nil
	}

	// Set port
	bc.port = port

	return &bc
}
