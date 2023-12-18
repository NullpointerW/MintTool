package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"os"
)

type Wallet struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

func main() {
	wei, ok := EthToWei("0.2")
	if ok {
		fmt.Println("0.2eth=", wei, "wei")
	}
}

func GenWallet(n int) {
	var wallets []Wallet
	wjr, err := os.ReadFile("wallet.json")
	if os.IsNotExist(err) {
		fmt.Println("no wallet exist,skip")
	} else if err != nil {
		fmt.Println("load wallet err:", err)
	} else {
		_ = json.Unmarshal(wjr, &wallets)
	}
	for i := 0; i < n; i++ {
		wallets = append(wallets, genWallet())
	}
	wjf, _ := os.Create("wallet.json")
	raw, _ := json.MarshalIndent(wallets, "", "    ")
	_, _ = wjf.Write(raw)
	_ = wjf.Close()
}

func genWallet() Wallet {
	// 生成一个新的私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// 获取私钥的字节形式
	privateKeyBytes := crypto.FromECDSA(privateKey)

	// 将私钥字节转换为十六进制字符串
	privateKeyHex := fmt.Sprintf("0x%x", privateKeyBytes)
	fmt.Println("Private Key:", privateKeyHex)

	// 从私钥生成公钥
	publicKey := privateKey.Public()

	// 从公钥生成公钥的字节形式
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	// 将公钥字节转换为十六进制字符串
	publicKeyHex := fmt.Sprintf("0x%x", publicKeyBytes[1:]) // 跳过ECDSA公钥前的0x04
	fmt.Println("Public Key:", publicKeyHex)

	// 从公钥生成以太坊地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("Address:", address)
	return Wallet{Address: address, PublicKey: publicKeyHex, PrivateKey: privateKeyHex}
}
