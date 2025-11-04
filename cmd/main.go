package main

import (
	"fmt"
	"log"
	"os"

	"Abby/api"
	"Abby/contracts"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file from any location")
	}

	var contractAddress []byte
	// 讀取合約地址
	if data, err := os.ReadFile("contract_address.txt"); err == nil {
		contractAddress = data
	} else {
		log.Fatal("Error reading contract address from any location")
	}

	// 連接到 Sepolia 測試網
	infuraURL := fmt.Sprintf("https://sepolia.infura.io/v3/%s", os.Getenv("INFURA_API_KEY"))
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 創建合約交互器
	interactor, err := contracts.NewContractInteractor(
		client,
		string(contractAddress),
		os.Getenv("PRIVATE_KEY"),
	)
	if err != nil {
		log.Fatal("Failed to create contract interactor:", err)
	}

	// 創建 API handler
	handler := api.NewStorageHandler(interactor)

	// 設置路由
	router := api.SetupRouter(handler)

	// 啟動服務器
	fmt.Println("Server is running on http://localhost:8081")
	fmt.Println("Swagger UI is available at http://localhost:8081/swagger/index.html")
	if err := router.Run(":8081"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

type NumArray struct {
	arr []int
}
