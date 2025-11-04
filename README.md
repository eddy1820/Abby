# 🧱 Abby Project

Abby is a sample smart contract implemented using Golang and Solidity.

---

## 🚀 Setup Project
Add a `.env` file with the following keys: `INFURA_API_KEY`, `PRIVATE_KEY`, and `SEPOLIA_RPC_URL`.


### 1️⃣ Generate swagger doc
```bash
swag init -g api/server.go
```
### 2️⃣ Start the server
```bash
go run cmd/main.go
```

🧩 Notes
Make sure you have Swagger installed before running the swag init command.
The server will run using the configuration specified in your .env file.
