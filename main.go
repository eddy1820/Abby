package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"Abby/contracts"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 檢查是否為預覽模式
	previewMode := false // 設置為 false 進行實際部署

	apiKey := os.Getenv("INFURA_API_KEY")
	fmt.Println("API Key: " + apiKey)

	// 連接到 Sepolia 測試網
	infuraURL := fmt.Sprintf("https://sepolia.infura.io/v3/%s", os.Getenv("INFURA_API_KEY"))
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 確認連接
	block, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current block number: %d\n", block)

	// 部署合約
	privateKey := os.Getenv("PRIVATE_KEY")

	if previewMode {
		fmt.Println("=== 預覽模式 ===")
		_, err := contracts.EstimateDeployment(client, privateKey)
		if err != nil {
			log.Fatal("Failed to estimate deployment:", err)
		}
		fmt.Println("要實際部署合約，請將 previewMode 設為 false")
	} else {
		fmt.Println("=== 部署模式 ===")
		_, err := contracts.DeployContract(client, privateKey)
		if err != nil {
			log.Fatal("Failed to deploy contract:", err)
		}
		fmt.Println("合約部署成功！")
		fmt.Println("Contract deployed successfully!")
	}
}
