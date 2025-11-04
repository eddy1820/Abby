package contracts

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ContractInteractor 用於與合約交互的結構體
type ContractInteractor struct {
	client   *ethclient.Client
	contract *Contracts
	auth     *bind.TransactOpts
	address  common.Address
}

// NewContractInteractor 創建新的合約交互器
func NewContractInteractor(client *ethclient.Client, contractAddress string, privateKey string) (*ContractInteractor, error) {
	// 轉換合約地址
	address := common.HexToAddress(contractAddress)

	// 創建合約實例
	contract, err := NewContracts(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract instance: %v", err)
	}

	// 轉換私鑰
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert private key: %v", err)
	}

	// 獲取鏈ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain id: %v", err)
	}

	// 創建交易選項
	auth, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	return &ContractInteractor{
		client:   client,
		contract: contract,
		auth:     auth,
		address:  address,
	}, nil
}

// GetValue 讀取當前存儲的值
func (ci *ContractInteractor) GetValue() (*big.Int, error) {
	value, err := ci.contract.Get(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get value: %v", err)
	}
	return value, nil
}

// SetValue 設置新的值
func (ci *ContractInteractor) SetValue(value *big.Int) error {
	// 獲取 gas 價格
	gasPrice, err := ci.client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %v", err)
	}
	ci.auth.GasPrice = gasPrice

	// 發送交易
	tx, err := ci.contract.Set(ci.auth, value)
	if err != nil {
		return fmt.Errorf("failed to set value: %v", err)
	}

	// 打印交易哈希
	log.Printf("Transaction sent: %s", tx.Hash().Hex())
	log.Printf("Waiting for transaction to be mined...")

	// 等待交易被確認
	receipt, err := bind.WaitMined(context.Background(), ci.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction: %v", err)
	}

	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	log.Printf("Transaction confirmed in block %d", receipt.BlockNumber)
	return nil
}

// WatchEvents 監聽合約事件
func (ci *ContractInteractor) WatchEvents() error {
	// 創建事件過濾器
	logs, err := ci.contract.FilterDataStored(&bind.FilterOpts{
		Start: 0, // 從最新的區塊開始
	})
	if err != nil {
		return fmt.Errorf("failed to filter events: %v", err)
	}
	defer logs.Close()

	// 打印歷史事件
	for logs.Next() {
		log.Printf("Historical event - New value stored: %s", logs.Event.NewValue)
	}

	// 創建事件訂閱
	sink := make(chan *ContractsDataStored)
	sub, err := ci.contract.WatchDataStored(&bind.WatchOpts{}, sink)
	if err != nil {
		return fmt.Errorf("failed to watch events: %v", err)
	}
	defer sub.Unsubscribe()

	// 監聽新事件
	for {
		select {
		case err := <-sub.Err():
			return fmt.Errorf("event subscription error: %v", err)
		case event := <-sink:
			log.Printf("New event - Value stored: %s", event.NewValue)
		}
	}
}
