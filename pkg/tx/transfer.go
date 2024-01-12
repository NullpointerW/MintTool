package tx

import (
	"context"
	"crypto/ecdsa"
	"ethereum-wallet-tool/pkg/util"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func Transfer(pk string, to string, val string, data []byte, ec *ethclient.Client) (txHash string, err error) {
	privateKeyHex := pk
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
	wei, ok := util.ToWei(val)
	if !ok {
		fmt.Println("invalid eth value:", val)
		return "", err
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
	toAddress := common.HexToAddress(to)

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
		return "", err
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
	if len(data) > 0 {
		tx.Data = data
	}
	txn := types.NewTx(tx)
	// 创建交易
	signedTx, err := types.SignTx(txn, types.NewCancunSigner(tx.ChainID), privateKey)
	if err != nil {
		return "", err
	}
	err = ec.SendTransaction(context.Background(), signedTx)
	return txn.Hash().String(), err
}
