package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	WallNum int
	Val     string
)

var ethwt = &cobra.Command{
	Use:   "ethwt",
	Short: "ethwt is a Ethereum wallet command-line tool ",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var gen = &cobra.Command{
	Use:   "gen",
	Short: "Generate wallet",
	Run: func(cmd *cobra.Command, args []string) {
		GenWallet(WallNum)
	},
}

var tx = &cobra.Command{
	Use:   "tx",
	Short: "Transfer to wallets",
	Run: func(cmd *cobra.Command, args []string) {
		Transfer(Val)
	},
}

var pa = &cobra.Command{
	Use:   "pa",
	Short: "Print private key as json array",
	Run: func(cmd *cobra.Command, args []string) {
		PrintPkJsonArray()
	},
}

func init() {
	tx.Flags().StringVarP(&Val, "val", "v", "0", "Value amount of ETH for each transfer")
	gen.Flags().IntVarP(&WallNum, "walletNum", "n", 1, "The number of wallets generated")
	ethwt.AddCommand(tx, gen, pa)
}

func main() {
	if err := ethwt.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
