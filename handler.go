package main

 import (
 	"errors"
 )

type Request struct {
	RType 	string 	`json:rType`
	Amount 	int		`json:amount`
	ID		string	`json:id`
}

func NewRequest() *Request {
	return &Request {
		RType : "0",
		Amount : 0,
		ID : "0",
	}
}

func (r *Request) Transaction(amount int, w *Wallet) (error){
	switch(r.RType){
	case "DEPOSIT": {
		newSum := w.Amount + amount
		w.Update(newSum)
		return nil
	}

	case "WITHDRAW":{
		newSum := w.Amount - amount
		w.Update(newSum)
		return nil
	}

	default:{
		return errors.New("Transaction type not found!")
	}

	}
}