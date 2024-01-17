package wal

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

func Gen() (Wallet, error) {
	// 生成一个新的私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return Wallet{}, fmt.Errorf("failed to generate private key: %v", err)
	}
	// 获取私钥的字节形式
	privateKeyBytes := crypto.FromECDSA(privateKey)

	// 将私钥字节转换为十六进制字符串
	privateKeyHex := fmt.Sprintf("0x%x", privateKeyBytes)
	//fmt.Println("Private Key:", privateKeyHex)

	// 从私钥生成公钥
	publicKey := privateKey.Public()

	// 从公钥生成公钥的字节形式
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return Wallet{}, errors.New("genWallet: cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	// 将公钥字节转换为十六进制字符串
	publicKeyHex := fmt.Sprintf("0x%x", publicKeyBytes[1:]) // 跳过ECDSA公钥前的0x04
	//fmt.Println("Public Key:", publicKeyHex)

	// 从公钥生成以太坊地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	//fmt.Println("Address:", address)
	return Wallet{Address: address, PublicKey: publicKeyHex, PrivateKey: privateKeyHex}, nil
}
