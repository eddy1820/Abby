package contracts

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EstimateDeployment 估算部署合約的成本但不實際部署
func EstimateDeployment(client *ethclient.Client, privateKeyHex string) (*Contracts, error) {
	// 轉換私鑰
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to convert private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to get public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %v", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain id: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	// 使用預設的 gas 限制
	gasEstimate := uint64(300000) // 標準合約部署大約需要 200,000-250,000 gas
	fmt.Printf("Using standard gas estimate: %d\n", gasEstimate)
	auth.GasLimit = gasEstimate

	// 計算預估的部署成本
	gasCost := new(big.Float).Mul(
		new(big.Float).SetInt(gasPrice),
		new(big.Float).SetUint64(auth.GasLimit),
	)
	ethCost := new(big.Float).Quo(gasCost, new(big.Float).SetUint64(1e18))
	fmt.Printf("Estimated deployment cost: %f ETH\n", ethCost)

	// 檢查錢包餘額
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %v", err)
	}
	balanceEth := new(big.Float).Quo(new(big.Float).SetInt(balance), new(big.Float).SetUint64(1e18))
	fmt.Printf("Wallet balance: %f ETH\n", balanceEth)

	// 檢查餘額是否足夠
	if balance.Cmp(new(big.Int).Mul(gasPrice, big.NewInt(int64(auth.GasLimit)))) < 0 {
		return nil, fmt.Errorf("insufficient funds for deployment")
	}

	fmt.Println("=== 部署預覽 ===")
	fmt.Printf("From address: %s\n", fromAddress.Hex())
	fmt.Printf("Gas price: %s Wei\n", gasPrice.String())
	fmt.Printf("Gas limit: %d\n", auth.GasLimit)
	fmt.Printf("Nonce: %d\n", nonce)

	return nil, nil
}

func DeployContract(client *ethclient.Client, privateKeyHex string) (*Contracts, error) {
	// 轉換私鑰
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to convert private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to get public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %v", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain id: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	// 使用預設的 gas 限制
	gasEstimate := uint64(300000) // 標準合約部署大約需要 200,000-250,000 gas
	fmt.Printf("Using standard gas estimate: %d\n", gasEstimate)
	auth.GasLimit = gasEstimate

	// 計算預估的部署成本
	gasCost := new(big.Float).Mul(
		new(big.Float).SetInt(gasPrice),
		new(big.Float).SetUint64(auth.GasLimit),
	)
	ethCost := new(big.Float).Quo(gasCost, new(big.Float).SetUint64(1e18))
	fmt.Printf("Estimated deployment cost: %f ETH\n", ethCost)

	// 檢查錢包餘額
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %v", err)
	}
	balanceEth := new(big.Float).Quo(new(big.Float).SetInt(balance), new(big.Float).SetUint64(1e18))
	fmt.Printf("Wallet balance: %f ETH\n", balanceEth)

	// 檢查餘額是否足夠
	if balance.Cmp(new(big.Int).Mul(gasPrice, big.NewInt(int64(auth.GasLimit)))) < 0 {
		return nil, fmt.Errorf("insufficient funds for deployment")
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // 不發送 ETH
	auth.GasPrice = gasPrice

	// 部署合約
	address, tx, instance, err := DeployContracts(auth, client)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy contract: %v", err)
	}

	// 打印部署信息
	fmt.Println("\n=== 部署成功 ===")
	fmt.Printf("合約地址: %s\n", address.Hex())
	fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())

	// 保存合約地址到文件
	addressFile := "contract_address.txt"
	err = os.WriteFile(addressFile, []byte(address.Hex()), 0644)
	if err != nil {
		log.Printf("Warning: Failed to save contract address: %v", err)
	}

	return instance, nil
}
