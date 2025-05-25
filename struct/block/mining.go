package block

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func (bc *Blockchain) Mining() bool {

	bc.mux.Lock()
	defer bc.mux.Unlock()
	if len(bc.transactionPool) == 0 {
		return false
	}
	_, err := bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, "MINING REWARD", MINING_REWARD, nil, nil)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return false
	}
	previousHash := bc.LastBlock().GetHash()
	bc.CreateBlock(bc.transactionPool, previousHash)
	if err := bc.SaveBlockchain(); err != nil {
		log.Printf("ERROR: Failed to save blockchain after mining: %v", err)
	}
	log.Println("action=mining, status=success")
	for _, n := range bc.neighbors {
		fmt.Println("Send consensus to neigbour ", n)
		endpoint := fmt.Sprintf("%s/consensus", n)
		client := &http.Client{}
		req, _ := http.NewRequest("PUT", endpoint, nil)
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("ERROR: %v", err)
			return false
		}
		log.Printf("%v", resp)
	}
	return true
}

func (bc *Blockchain) StartMining() {
	bc.Mining()
	_ = time.AfterFunc(time.Second*MINING_TIMER_SEC, bc.StartMining)
}

// ValidProof validates the proof of work.
func (bc *Blockchain) ValidProof(nonce int, previousHash string, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := NewBlock(transactions, previousHash)
	guessBlock.SetNonce(nonce)
	guessHash := guessBlock.CalculateHash()
	return strings.HasPrefix(guessHash, zeros)
}

// ProofOfWork finds the proof of work.
func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().GetHash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) ValidChain(chain []*Block) bool {

	preBlock := chain[0]
	currentIndex := 1
	for currentIndex < len(chain) {
		b := chain[currentIndex]
		if b.GetPrevHash() != preBlock.GetHash() {
			return false
		}

		if !bc.ValidProof(b.GetNonce(), b.GetPrevHash(), b.GetTransactions(), MINING_DIFFICULTY) {
			return false
		}

		preBlock = b
		currentIndex += 1
	}
	return true
}

func (bc *Blockchain) RegisterNewWallet(blockchainAddress string, message string) bool {

	_, err := bc.AddTransaction(MINING_SENDER, blockchainAddress, message, 0, nil, nil)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return false
	}

	bc.StartMining()

	return true
}

func (bc *Blockchain) CalculateTotalBalance(blockchainAddress string) (float32, error) {
	var totalBalance float32 = 0.0
	addressFound := false

	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value

			if blockchainAddress == t.recipientBlockchainAddress {
				totalBalance += value
				addressFound = true
			}

			if blockchainAddress == t.senderBlockchainAddress {
				totalBalance -= value
				addressFound = true
			}
		}
	}

	if !addressFound {
		return 0.0, fmt.Errorf("Address not found in the Blockchain")
	}

	return totalBalance, nil
}
