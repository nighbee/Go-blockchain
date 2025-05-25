package wallet

import (
	"block/struct/utils"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

type Transaction struct {
	message                    string
	recipientBlockchainAddress string
	senderBlockchainAddress    string
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	value                      float32
}

type TransactionRequest struct {
	Message                    *string  `json:"message"`
	RecipientBlockchainAddress *string  `json:"recipientBlockchainAddress"`
	SenderBlockchainAddress    *string  `json:"senderBlockchainAddress"`
	SenderPublicKey            *string  `json:"senderPublicKey"`
	SenderPrivateKey           *string  `json:"senderPrivateKey"`
	Value                      *float32 `json:"value"`
}

func NewWallet() *Wallet {
	w := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey

	h2 := sha256.New()
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)

	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])

	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	chsum := digest6[:4]

	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum[:])

	address := base58.Encode(dc8)
	w.blockchainAddress = address
	return w
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey        string `json:"privateKey"`
		PublicKey         string `json:"publicKey"`
		BlockchainAddress string `json:"blockchainAddress"`
	}{
		PrivateKey:        w.PrivateKeyStr(),
		PublicKey:         w.PublicKeyStr(),
		BlockchainAddress: w.BlockchainAddress(),
	})
}

func NewTransaction(
	message string,
	sender string,
	recipient string,
	privateKey *ecdsa.PrivateKey,
	publicKey *ecdsa.PublicKey,
	value float32) *Transaction {
	return &Transaction{message, recipient, sender, privateKey, publicKey, value}
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(t)

	log.Println("Generate signature", string(m))

	h := sha256.Sum256([]byte(m))
	r, s, err := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	if err != nil {
		fmt.Println("Signing signature failed: ", err)
		return nil
	}
	return &utils.Signature{R: r, S: s}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message   string  `json:"message"`
		Sender    string  `json:"senderBlockchainAddress"`
		Recipient string  `json:"recipientBlockchainAddress"`
		Value     float32 `json:"value"`
	}{
		Message:   t.message,
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderPrivateKey == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Message == nil ||
		tr.Value == nil {
		return false
	}
	return true
}
