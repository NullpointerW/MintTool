package main

import (
	"encoding/json"
	"fmt"
	_tx "github.com/NullpointerW/ethereum-wallet-tool/pkg/tx"
	"github.com/ethereum/go-ethereum/ethclient"
	"os"
	"strings"
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
		_, err = _tx.WaitForTransactionConfirmation(ec, txHash)
		if err != nil {
			txR.ErrMsg = err.Error()
		}
	}
	return txR
}
