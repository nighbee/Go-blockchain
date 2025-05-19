# Blockchain Server

![Go](https://img.shields.io/badge/Go-1.20-blue)
![React](https://img.shields.io/badge/React-18.2.0-blue)
![License](https://img.shields.io/badge/License-MIT-green)

## �� Description

Blockchain Server is a blockchain implementation in Go with a React frontend. The project supports block creation, mining, wallet registration, and transaction management through REST API. A key feature is the ability to save blockchain state to a JSON file and reset the blockchain to its initial state.

---

## �� Key Features

- 🔗 **Block Mining** — creation of new blocks with transactions and mining rewards
- 👛 **Wallet Management** — registration and tracking of wallets in the blockchain
- �� **Transactions** — creation and validation of transactions between wallets
- 🔄 **State Persistence** — automatic saving of blockchain state to JSON file
- 🗑️ **Blockchain Reset** — ability to reset to initial state
- 🤝 **Consensus** — conflict resolution between network nodes
- 🌐 **REST API** — comprehensive API for blockchain interaction
- �� **Security** — ECDSA signatures for transaction validation
- 📊 **Real-time Updates** — live blockchain state monitoring

---

## 🛠️ Technologies

- **Backend**: Go
- **Frontend**: React
- **HTTP Framework**: Gorilla Mux
- **Data Storage**: JSON files
- **Cryptography**: ECDSA for transaction signing
- **API**: RESTful architecture

---

## 📂 Project Structure

```plaintext
├── struct/
│   ├── block/             # Blockchain logic
│   │   ├── blockchain.go  # Core blockchain logic
│   │   ├── mining.go      # Block mining
│   │   ├── transaction.go # Transactions
│   │   └── storage.go     # State persistence
│   ├── wallet/            # Wallet implementation
│   └── utils/             # Utility functions
├── server/
│   ├── handlers/          # API handlers
│   └── middleware/        # CORS and logging
├── src/                   # Frontend React application
├── data/                  # Blockchain data storage
├── main.go               # Entry point
└── go.mod                # Go dependencies
```

---

## 📦 Installation and Setup

### Backend

1. Install Go
2. Clone the repository:
```bash
git clone https://github.com/your-username/blockchain-server.git
cd blockchain-server
```

3. Install dependencies:
```bash
go mod tidy
```

4. Start the server:
```bash
go run main.go
```

### Frontend

1. Navigate to src directory:
```bash
cd src
```

2. Install dependencies:
```bash
npm install
```

3. Start the application:
```bash
npm start
```

---

## 📖 API Documentation

### REST API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /chain | Get current blockchain |
| GET | /wallets | Get list of registered wallets |
| POST | /wallet/register | Register new wallet |
| POST | /transactions | Add new transaction |
| GET | /balance | Get wallet balance |
| POST | /mine | Mine new block |
| POST | /reset | Reset blockchain to initial state |
| GET | /consensus | Run consensus process |
| GET | /nodes | Get connected nodes |
| POST | /sign | Sign transaction |

### API Examples

Register a wallet:
```bash
curl -X POST http://localhost:5001/wallet/register \
-H "Content-Type: application/json"
```

Create a transaction:
```bash
curl -X POST http://localhost:5001/transactions \
-H "Content-Type: application/json" \
-d '{
  "sender": "sender_address",
  "recipient": "recipient_address",
  "value": 10.0,
  "message": "Payment"
}'
```

Reset blockchain:
```bash
curl -X POST http://localhost:5001/reset
```

Get wallet balance:
```bash
curl -X GET http://localhost:5001/balance?address=wallet_address
```

Mine a block:
```bash
curl -X POST http://localhost:5001/mine \
-H "Content-Type: application/json" \
-d '{
  "minerAddress": "miner_address"
}'
```

---

## 🧪 Testing

1. Start the server
2. Open web interface at http://localhost:3000
3. Test core functionality:
   - Wallet registration
   - Transaction creation
   - Block mining
   - Blockchain reset
   - API endpoints using curl or Postman

---

## 🤝 Contributing

We welcome contributions!
- Create issues to discuss problems
- Submit pull requests with improvements
- Share ideas for project development

---

## 📜 License

This project is licensed under the MIT License.

---

## 📧 Contact

Author: Almaz Toktassin
- 📬 Email: almaztok8@gmail.com
- 💻 GitHub: nighbee
