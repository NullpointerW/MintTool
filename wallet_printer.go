package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func PrintPkJsonArray() {
	var wallets []Wallet
	wjr, err := os.ReadFile("wallet.json")
	if os.IsNotExist(err) {
		fmt.Println("no wallet found")
		return
	} else if err != nil {
		fmt.Println("load wallet err:", err)
		return
	} else {
		_ = json.Unmarshal(wjr, &wallets)
	}
	var jsArray []string
	for _, w := range wallets {
		jsArray = append(jsArray, strings.TrimPrefix(w.PrivateKey, "0x"))
	}
	pkjf, _ := os.Create("pks.json")
	raw, _ := json.MarshalIndent(jsArray, "", "    ")
	_, _ = pkjf.Write(raw)
	_ = pkjf.Close()
}
