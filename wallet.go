package main

import (
	"errors"
)

type Wallet struct{
	WalletId 	string 	`json:"walletid"`
	Amount 		int 	`json:"amount"`
}

type Storage struct{
	WalletStorage map[string]Wallet
}

func NewWallet(newID string, newAmount int) Wallet{
	return Wallet {
		WalletId : newID,
		Amount	 : newAmount,
	}
}

func NewStorage() *Storage {
	return &Storage {
		make(map[string]Wallet),
	}
}

func (s *Storage) FillStorage() {
	s.WalletStorage["a"] = NewWallet("a", 10)
	s.WalletStorage["b"] = NewWallet("b", 20)
	s.WalletStorage["c"] = NewWallet("c", 30)
}

func (s *Storage) Get(id string) (Wallet, error){
	for i := range s.WalletStorage {
		if id == s.WalletStorage[i].WalletId {
			return s.WalletStorage[i], nil
		}
	}
	return NewWallet("", 0), errors.New("Wallet not found")
}

func (w *Wallet) Update(newAmount int) (*Wallet, error){
		w.Amount = newAmount
		return w, nil
}
