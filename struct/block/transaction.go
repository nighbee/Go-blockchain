package block

import (
	"block/struct/utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Transaction struct {
	message                    string
	recipientBlockchainAddress string
	senderBlockchainAddress    string
	value                      float32
}

type TransactionRequest struct {
	Message                    *string  `json:"message"`
	RecipientBlockchainAddress *string  `json:"recipientBlockchainAddress"`
	SenderBlockchainAddress    *string  `json:"senderBlockchainAddress"`
	SenderPublicKey            *string  `json:"senderPublicKey"`
	Signature                  *string  `json:"signature"`
	Value                      *float32 `json:"value"`
}

type BalanceResponse struct {
	Balance float32 `json:"balance"`
	Error   string  `json:"error"`
}

func (bc *Blockchain) AddTransaction(sender string,
	recipient string,
	message string,
	value float32,
	senderPublicKey *ecdsa.PublicKey,
	s *utils.Signature) (bool, error) {

	t := &Transaction{
		senderBlockchainAddress:    sender,
		recipientBlockchainAddress: recipient,
		message:                    message,
		value:                      value,
	}

	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true, nil
	}

	if !bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		return false, fmt.Errorf("ERROR: Verify Transaction")
	}

	balance, err := bc.CalculateTotalBalance(sender)
	if err != nil {
		return false, fmt.Errorf("ERROR: CalculateTotalAmount: %v", err)
	}

	if balance < value {
		return false, fmt.Errorf("ERROR: Not enough balance in a wallet")
	}

	bc.transactionPool = append(bc.transactionPool, t)

	if err := bc.SaveBlockchain(); err != nil {
		log.Printf("ERROR: Failed to save blockchain after adding transaction: %v", err)
	}
	return true, nil
}

func (bc *Blockchain) ClearTransactionPool() {
	bc.transactionPool = bc.transactionPool[:0]
}

func (bc *Blockchain) CreateTransaction(sender string, recipient string, message string, value float32,
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature) (bool, error) {

	isTransacted, err := bc.AddTransaction(sender, recipient, message, value, senderPublicKey, s)

	if err != nil {

		log.Printf("ERROR: %v", err)
		return false, err
	}

	if isTransacted {
		for _, n := range bc.neighbors {
			publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(),
				senderPublicKey.Y.Bytes())
			signatureStr := s.String()
			bt := &TransactionRequest{
				&message, &publicKeyStr, &recipient, &sender, &signatureStr, &value}
			m, _ := json.Marshal(bt)
			buf := bytes.NewBuffer(m)
			endpoint := fmt.Sprintf("%s/transactions", n)
			client := &http.Client{}
			req, _ := http.NewRequest("PUT", endpoint, buf)
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("ERROR: %v", err)
				return false, err
			}
			log.Printf("%v", resp)
		}
	}

	return isTransacted, nil
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(
				t.senderBlockchainAddress,
				t.recipientBlockchainAddress,
				t.message,
				t.value))
	}
	return transactions
}

func NewTransaction(sender string, recipient string, message string, value float32) *Transaction {
	return &Transaction{
		senderBlockchainAddress:    sender,
		recipientBlockchainAddress: recipient,
		message:                    message,
		value:                      value,
	}
}

func (bc *Blockchain) TransactionPool() []*Transaction {
	return bc.transactionPool
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Message == nil ||
		tr.Value == nil ||
		tr.Signature == nil {
		return false
	}
	return true
}
func (bc *Blockchain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, err := json.Marshal(t)
	if err != nil {
		log.Printf("ERROR: Failed to marshal transaction for verification: %v", err)
		return false
	}

	log.Printf("Verifying transaction data: %s", string(m))

	hash := sha256.Sum256(m)
	log.Printf("Verification hash: %x", hash)

	log.Printf("Verifying signature R: %x", s.R)
	log.Printf("Verifying signature S: %x", s.S)

	verified := ecdsa.Verify(senderPublicKey, hash[:], s.R, s.S)
	if !verified {
		log.Printf("ERROR: Signature verification failed for publicKey: %x, signature: %s", senderPublicKey, s.String())
		log.Printf("Public key components - X: %x, Y: %x", senderPublicKey.X, senderPublicKey.Y)
	}
	return verified
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" senderBlockchainAddress      %s\n", t.senderBlockchainAddress)
	fmt.Printf(" recipientBlockchainAddress   %s\n", t.recipientBlockchainAddress)
	fmt.Printf(" message                      %s\n", t.message)
	fmt.Printf(" value                          %.1f\n", t.value)
}

func (br *BalanceResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Balance float32 `json:"balance"`
		Error   string  `json:"error"`
	}{
		Balance: br.Balance,
		Error:   br.Error,
	})
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message   string  `json:"message"`
		Recipient string  `json:"recipientBlockchainAddress"`
		Sender    string  `json:"senderBlockchainAddress"`
		Value     float32 `json:"value"`
	}{
		Message:   t.message,
		Recipient: t.recipientBlockchainAddress,
		Sender:    t.senderBlockchainAddress,
		Value:     t.value,
	})
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	v := &struct {
		Message   *string  `json:"message"`
		Recipient *string  `json:"recipientBlockchainAddress"`
		Sender    *string  `json:"senderBlockchainAddress"`
		Value     *float32 `json:"value"`
	}{
		Message:   &t.message,
		Recipient: &t.recipientBlockchainAddress,
		Sender:    &t.senderBlockchainAddress,
		Value:     &t.value,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}
