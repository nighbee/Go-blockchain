# 🧱 Go-Blockchain

A fully functional blockchain system written in **Go (Golang)** with a **React-based web interface**. This project includes core blockchain features such as Proof of Work, wallet management with ECDSA keys, transaction processing, a REST API, and multi-node network simulation.

---

## 🚀 Features

### ✅ Core Blockchain Features
- Block creation with SHA-256 hashing
- Mining with Proof-of-Work (configurable difficulty)
- Wallets with public/private key pairs (ECDSA)
- Secure transaction processing with digital signatures
- Persistent blockchain state (JSON-based storage)

### 🌐 Web Interface (Frontend)
- Blockchain Explorer (View all blocks and transactions)
- Wallet creation and management
- Real-time transaction creation and monitoring
- Mining interface
- Node overview panel

### 🌍 Network Features
- Multi-node simulation (localhost ports)
- Neighbor discovery and dynamic sync
- Consensus algorithm to resolve forks
- RESTful API to interact with blockchain

---

## 📂 Project Structure
```bash
Go-blockchain/
├── backend/
│ ├── blockchain/ # Core blockchain logic (blocks, transactions, mining, wallets)
│ ├── server/ # API server and handlers
│ ├── middleware/ # CORS and logging
│ └── main.go # Application entry point
│
├── frontend/
│ ├── src/components/ # React components (Wallet, Transactions, Blockchain, Mining)
│ ├── src/api.js # API client
│ └── src/App.js # Main app routing
│
├── data/ # Blockchain state (JSON file)
└── README.md
```

---

## ⚙️ Installation

### 🔧 Prerequisites
- [Go](https://go.dev/doc/install)
- [Node.js & npm](https://nodejs.org/)
- Git

### 📦 Backend Setup

```bash
git clone https://github.com/nighbee/Go-blockchain.git
cd Go-blockchain/backend
```
# Install dependencies
go mod tidy

# Run node on port 5001
```bash
PORT=5001 go run main.go
```
To simulate multiple nodes (on ports 5002 and 5003), use the provided batch script (run_nodes.bat) or manually set different PORT values.

💻 Frontend Setup
```bash
cd Go-blockchain/frontend
npm install
npm start
```
The frontend will be available at http://localhost:3000.

🧪 API Endpoints (Sample)
```bash
Endpoint	        Method	Description
/chain	          GET	    Get full blockchain
/transactions	    POST	  Submit a new transaction
/wallet/register	POST	  Create a new wallet
/mine	            GET	    Mine a new block
/balance?address=	GET	    Get wallet balance
/consensus	      PUT	    Trigger consensus among nodes
/reset	          POST	  Reset blockchain to genesis state
```
🛡️ Security Features
ECDSA digital signatures for transaction authentication

Hash validation for block integrity

Balance verification before approving transactions

Mutex locks to prevent race conditions in concurrent mining or syncing

🧠 Key Concepts Learned
Blockchain architecture and decentralization

Go concurrency (mutex, goroutines)

REST API design

Digital signature cryptography (ECDSA)

React component-based UI design
