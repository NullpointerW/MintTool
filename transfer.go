package main

import (
	"context"
	"encoding/json"
	"fmt"
	_tx "github.com/NullpointerW/ethereum-wallet-tool/pkg/tx"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"os"
	"strings"
	"time"
)

type TxRecord struct {
	Hash    string `json:"hash"`
	Address string `json:"address"`
	PK      string `json:"PK"`
	Value   string `json:"value"`
	//GasFee  string `json:"gasFee"`
	ErrMsg string `json:"errMsg"`
}

func Transfer(val string) {
	ec, err := ethclient.Dial(Endpoint)
	if err != nil {
		fmt.Printf("Failed to connect to the Ethereum client: %v", err)
	}
	var wallets []Wallet
	wjr, err := os.ReadFile("wallet.json")
	if os.IsNotExist(err) {
		fmt.Println("no wallet found")
		return
	} else if err != nil {
		fmt.Println("load wallet err:", err)
		return
	}
	_ = json.Unmarshal(wjr, &wallets)
	pkb, err := os.ReadFile(".PK")
	if err != nil {
		fmt.Println("load wallet err:", err)
		return
	}
	pk := strings.TrimPrefix(string(pkb), "0x")
	var txRecords []TxRecord
	for _, w := range wallets {
		record := transfer(pk, w, val, ec)
		//if err != nil {
		//	fmt.Println("transfer err:", err)
		//	continue
		//}
		fmt.Printf("%#+v\n", record)
		txRecords = append(txRecords, record)
	}
	if len(txRecords) > 0 {
		txjf, _ := os.Create("txRecords.json")
		raw, _ := json.MarshalIndent(txRecords, "", "    ")
		_, _ = txjf.Write(raw)
		_ = txjf.Close()
	}

}

func transfer(mpk string, w Wallet, val string, ec *ethclient.Client) TxRecord {
	txHash, err := _tx.Transfer(mpk, w.Address, val, nil, ec)
	txR := TxRecord{
		Hash:    txHash.String(),
		Address: w.Address,
		PK:      w.PrivateKey,
		Value:   val,
	}
	if err != nil {
		txR.ErrMsg = err.Error()
	} else {
		_, err = waitForTransactionConfirmation(ec, txHash)
		if err != nil {
			txR.ErrMsg = err.Error()
		}
	}
	return txR
}
func waitForTransactionConfirmation(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	var receipt *types.Receipt
	var err error

	// 设置查询上下文，可以设置超时
	ctx := context.Background()

	// 设置轮询间隔

	fmt.Println("waiting for tx confirmed")
	for {
		receipt, err = client.TransactionReceipt(ctx, txHash)
		if err != nil {
			if err == ethereum.NotFound {
				// 如果收据未找到，继续轮询
				time.Sleep(time.Second)
				continue
			}
			// 如果发生其他错误，则返回错误
			return nil, err
		}
		// 如果收据不为空，表示交易已被确认
		if receipt != nil {
			break
		}
	}

	return receipt, nil
}
