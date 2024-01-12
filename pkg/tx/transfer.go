package tx

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/NullpointerW/ethereum-wallet-tool/pkg/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

func Transfer(pk string, to string, val string, data []byte, ec *ethclient.Client) (txHash common.Hash, err error) {
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
		return common.Hash{}, err
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
		Data:     data,
	}

	estimatedGas, err := ec.EstimateGas(context.Background(), msg)
	if err != nil {
		return common.Hash{}, err
	}
	fmt.Println("maxPriorityFeePerGas:", gasTipCap, "GasFeeCap:", maxFeePerGas, "Gas:", estimatedGas)
	cid, err := ec.ChainID(context.Background())
	if err != nil {
		return common.Hash{}, fmt.Errorf("get chainID error:%w", err)
	}
	tx := &types.DynamicFeeTx{
		ChainID:   cid,
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
		return common.Hash{}, err
	}
	//fmt.Println("signedTx.Hash()", signedTx.Hash().String())
	//fmt.Println("txn.Hash()", txn.Hash().String())
	err = ec.SendTransaction(context.Background(), signedTx)
	return signedTx.Hash(), err
}
func WaitForTransactionConfirmation(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	var receipt *types.Receipt
	var err error

	// 设置查询上下文，可以设置超时
	ctx := context.Background()

	// 设置轮询间隔

	for {
		receipt, err = client.TransactionReceipt(ctx, txHash)
		if err != nil {
			if err == ethereum.NotFound {
				fmt.Println("waiting for confirmed tx:", txHash.String())
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
