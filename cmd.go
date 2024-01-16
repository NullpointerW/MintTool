package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var (
	Endpoint string
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
		wallNum := 0
		if len(args) > 0 {
			var err error
			wallNum, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("generate wallet failed:", err)
				return
			}
		}
		GenWallet(wallNum)
	},
	Args: cobra.MinimumNArgs(1),
}

var tx = &cobra.Command{
	Use:   "tx",
	Short: "Transfer to wallets",
	Run: func(cmd *cobra.Command, args []string) {
		val := "0"
		if len(args) > 0 {
			val = args[0]
		}
		Transfer(val)
	},
	Args: cobra.MinimumNArgs(1),
}

var pa = &cobra.Command{
	Use:   "pa",
	Short: "Print private key as json array",
	Run: func(cmd *cobra.Command, args []string) {
		PrintPkJsonArray()
	},
}
var pl = &cobra.Command{
	Use:   "pl",
	Short: "Print private key as line txt",
	Run: func(cmd *cobra.Command, args []string) {
		PrintPkLine()
	},
}

func init() {
	tx.Flags().StringVarP(&Endpoint, "endpoint", "e", "https://ethereum.publicnode.com", "Network rpc endpoint")
	ethwt.AddCommand(tx, gen, pa, pl)
}

func main() {
	if err := ethwt.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
