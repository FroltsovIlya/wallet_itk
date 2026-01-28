package main

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var walletCache sync.Map //cache storage for stabilisation requests

type Handler struct {		
	storage StorageInterface //interface because we have tests
}

func NewHandler(s StorageInterface) *Handler {
	return &Handler{storage: s}	//classic constructor
}

type Request struct {	//structure of our request
	WalletID      string `json:"walletId"`
	OperationType string `json:"operationType"`
	Amount        int64  `json:"amount"`
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/v1/wallet", h.updateWallet) //realises wallet update
	r.GET("/api/v1/wallets/:id", h.getWallet)//getting wallet from db or cache
}

func (h *Handler) updateWallet(c *gin.Context) { //func to update wallet xd
	var req Request
	if err := c.BindJSON(&req); err != nil { //creating request and parsing from json
		c.JSON(400, gin.H{"error": "bad request"}) //if error of parsing, bad request
		return
	}

	if req.Amount <= 0 { //if amount of transaction negative or 0, error
		c.JSON(400, gin.H{"error": "amount must be positive"})
		return
	}

	var wallet *Wallet //create wallet struct for handle
	var err error //also error

	if req.OperationType == "DEPOSIT" { //chose operation type
		wallet, err = h.storage.Deposit(req.WalletID, req.Amount)
	} else if req.OperationType == "WITHDRAW" {
		wallet, err = h.storage.Withdraw(req.WalletID, req.Amount)
	} else {
		c.JSON(400, gin.H{"error": "unknown operation"}) //if request type unknown, error
		return
	}

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	walletCache.Store(wallet.ID, wallet) //update cache

	c.JSON(200, wallet) //show updated wallet
}

func (h *Handler) getWallet(c *gin.Context) { //getting wallet, i love it
	id := c.Param("id") //getting id from context

	if v, ok := walletCache.Load(id); ok { //load wallet from cache
		c.JSON(200, v) //if exists, return from cache
		return
	}

	wallet, err := h.storage.Get(id) //if in cache not found, get from db
	if err != nil { //got error
		if err.Error() == "wallet not found" { //if error of not existed wallet
			c.JSON(404, gin.H{"error": "wallet not found"}) //return 404
			return
		}
		c.JSON(500, gin.H{"error": "internal error"}) //else this error in db may be
		return
	}

	walletCache.Store(id, wallet) //add wallet to cache
	c.JSON(200, wallet) //return wallet
}
