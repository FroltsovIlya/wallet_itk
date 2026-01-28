package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"

	"github.com/gin-gonic/gin"
)
 //all these tests like for in Get()
func TestRequest_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	storage := newFakeStorage()
	handler := NewHandler(storage)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/wallet",
		bytes.NewBuffer([]byte(`{bad json`)),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	} else {
		t.Log("bad json test ok")
	}
}


func TestDeposit_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	storage := newFakeStorage()
	handler := NewHandler(storage)
	handler.RegisterRoutes(router)

	body := `{
		"walletId": "a",
		"operationType": "DEPOSIT",
		"amount": 500
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/wallet",
		bytes.NewBufferString(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if storage.wallets["a"].Amount != 1500 {
		t.Fatalf("expected balance 1500, got %d", storage.wallets["a"].Amount)
	} else {
		t.Log("deposit ok")
	}
}


func TestDeposit_BadAmount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	storage := newFakeStorage()
	handler := NewHandler(storage)
	handler.RegisterRoutes(router)

	body := `{
		"walletId": "a",
		"operationType": "DEPOSIT",
		"amount": -100
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/wallet",
		bytes.NewBufferString(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	} else {
		t.Log("bad request status. ok")
	}
}


func TestWithdraw_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	storage := newFakeStorage()
	handler := NewHandler(storage)
	handler.RegisterRoutes(router)

	body := `{
		"walletId": "a",
		"operationType": "WITHDRAW",
		"amount": 300
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/wallet",
		bytes.NewBufferString(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if storage.wallets["a"].Amount != 700 {
		t.Fatalf("expected balance 700, got %d", storage.wallets["a"].Amount)
	} else {
		t.Log("withdraw ok")
	}
}


func TestWithdraw_BadAmount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	storage := newFakeStorage()
	handler := NewHandler(storage)
	handler.RegisterRoutes(router)

	body := `{
		"walletId": "a",
		"operationType": "WITHDRAW",
		"amount": 5000
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/wallet",
		bytes.NewBufferString(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	} else {
		t.Log("staus 400. ok")
	}
}
