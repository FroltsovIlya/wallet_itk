package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeStorage struct {
	wallets map[string]*Wallet
}

func newFakeStorage() *fakeStorage {
	return &fakeStorage{
		wallets: map[string]*Wallet{
			"a": {ID: "a", Amount: 1000},
		},
	}
}

func (f *fakeStorage) Get(id string) (*Wallet, error) {
	w, ok := f.wallets[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return w, nil
}

func (f *fakeStorage) Deposit(id string, amount int64) (*Wallet, error) {
	w, ok := f.wallets[id]
	if !ok {
		return nil, errors.New("not found")
	}
	w.Amount += amount
	return w, nil
}

func (f *fakeStorage) Withdraw(id string, amount int64) (*Wallet, error) {
	w, ok := f.wallets[id]
	if !ok {
		return nil, errors.New("not found")
	}
	if w.Amount < amount {
		return nil, errors.New("not enough money")
	}
	w.Amount -= amount
	return w, nil
}


func TestGetWallet_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	storage := newFakeStorage()
	handler := NewHandler(storage)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/a", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	} else {
		t.Log("got 200 how planned")
	}
}

func TestGetWallet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	storage := newFakeStorage()
	handler := NewHandler(storage)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/unknown", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	} else {
		t.Log("got 404 how planned")
	}
}


