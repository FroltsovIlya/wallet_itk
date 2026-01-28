package main

type Wallet struct { //struct of wallet, nothing special
	ID     string `json:"walletId"`
	Amount int64  `json:"amount"`
}
