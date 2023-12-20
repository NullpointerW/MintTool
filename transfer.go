package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
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
	ec, err := ethclient.Dial("https://ethereum.publicnode.com")
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
		record, err := transfer(pk, w, val, ec)
		if err != nil {
			fmt.Println("transfer err:", err)
			continue
		}
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

func transfer(mpk string, w Wallet, val string, ec *ethclient.Client) (TxRecord, error) {
	privateKeyHex := mpk
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	// 获取账户地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取账户的nonce
	nonce, err := ec.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get account nonce: %v", err)
	}

	// 设置转账金额（0.001 ETH）
	wei, ok := EthToWei(val)
	if !ok {
		fmt.Println("invalid eth value:", val)
		return TxRecord{}, errors.New(fmt.Sprintf("invalid eth value: %s", val))
	}

	gasTipCap, err := ec.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas tip cap: %v", err)
	}

	// 设置 maxFeePerGas。这通常是 baseFeePerGas + maxPriorityFeePerGas，但要留有余地以适应 baseFee 的变动
	baseFee, err := ec.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to get latest block header: %v", err)
	}
	maxFeePerGas := new(big.Int).Add(baseFee.BaseFee, gasTipCap)

	// 设置接收方地址
	toAddress := common.HexToAddress(w.Address)

	// 估算 Gas 量
	msg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &toAddress,
		Gas:      0,
		GasPrice: maxFeePerGas,
		Value:    wei,
		Data:     nil,
	}

	estimatedGas, err := ec.EstimateGas(context.Background(), msg)
	if err != nil {
		return TxRecord{}, fmt.Errorf("failed to estimate gas: %v", err)
	}
	fmt.Println("maxPriorityFeePerGas:", gasTipCap, "GasFeeCap:", maxFeePerGas, "Gas:", estimatedGas)
	tx := &types.DynamicFeeTx{
		ChainID:   big.NewInt(1),
		Nonce:     nonce,
		GasTipCap: gasTipCap,    // maxPriorityFeePerGas
		GasFeeCap: maxFeePerGas, // max Fee
		Gas:       estimatedGas,
		To:        &toAddress,
		Value:     wei,
	}

	// 创建交易
	signedTx, err := types.SignTx(types.NewTx(tx), types.NewCancunSigner(tx.ChainID), privateKey)
	if err != nil {
		return TxRecord{}, errors.New(fmt.Sprintf("signTx err:%s", err.Error()))
	}
	err = ec.SendTransaction(context.Background(), signedTx)
	txR := TxRecord{
		Hash:    signedTx.Hash().String(),
		Address: w.Address,
		PK:      w.PrivateKey,
		Value:   val,
	}
	if err != nil {
		txR.ErrMsg = err.Error()
	}
	return txR, nil
}
