ethereum-wallet-tool
-----------------------
## function
* Bulk Wallet Generation
* Batch transfer to wallet
## install
```shell
go build ewt.exe/ewt
```
## usage
Command help:
```shell
ewt.exe --help
```
1. Configure the primary wallet key in the `.PK` file.
2. Generate a specified number of wallets in bulk,and Wallet information will be logged as JSON.
```shell
ewt.exe gen {num}
```
3. Use the main wallet to batch transfer a specified amount of eth to all wallets in wallet.json,and txRecord record will be logged as JSON.
```shell
ewt.exe tx {value}
```
