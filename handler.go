package main

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage *Storage
}

func NewHandler(s *Storage) *Handler {
	return &Handler{storage: s}
}

type Request struct {
	WalletID      string `json:"walletId"`
	OperationType string `json:"operationType"`
	Amount        int64  `json:"amount"`
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/v1/wallet", h.updateWallet)
	r.GET("/api/v1/wallets/:id", h.getWallet)
}

func (h *Handler) updateWallet(c *gin.Context) {
	var req Request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}

	if req.Amount <= 0 {
		c.JSON(400, gin.H{"error": "amount must be positive"})
		return
	}

	var wallet *Wallet
	var err    error

	if req.OperationType == "DEPOSIT" {
		wallet, err = h.storage.Deposit(req.WalletID, req.Amount)
	} else if req.OperationType == "WITHDRAW" {
		wallet, err = h.storage.Withdraw(req.WalletID, req.Amount)
	} else {
		c.JSON(400, gin.H{"error": "unknown operation"})
		return
	}

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, wallet)
}

func (h *Handler) getWallet(c *gin.Context) {
	id := c.Param("id")

	wallet, err := h.storage.Get(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "wallet not found"})
		return
	}

	c.JSON(200, wallet)
}
