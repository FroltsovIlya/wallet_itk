package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"log"
	"time"
)

type Storage struct {
	db *sql.DB
}

type StorageInterface interface {
	Get(id string) (*Wallet, error)
	Deposit(id string, amount int64) (*Wallet, error)
	Withdraw(id string, amount int64) (*Wallet, error)
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) InitDB() {
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

    var err error
    s.db, err = sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }

    err = s.db.Ping()
    if err != nil {
        log.Println("waiting for database...")
        time.Sleep(5 * time.Second)
        }

    fmt.Println("Successfully connected to PostgreSQL")

    _, err = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS wallets (
            id TEXT PRIMARY KEY,
            balance BIGINT NOT NULL DEFAULT 0
        )
    `)
    if err != nil {
        log.Fatal("failed to create wallets table:", err)
    }

    _, err = s.db.Exec(`
        INSERT INTO wallets (id, balance)
        VALUES ('a', 1000)
        ON CONFLICT (id) DO NOTHING
    `)
    if err != nil {
        log.Fatal("failed to insert test wallet:", err)
    }
	_, err = s.db.Exec(`
        INSERT INTO wallets (id, balance)
        VALUES ('b', 2000)
        ON CONFLICT (id) DO NOTHING
    `)
    if err != nil {
        log.Fatal("failed to insert test wallet:", err)
    }
	_, err = s.db.Exec(`
        INSERT INTO wallets (id, balance)
        VALUES ('c', 3000)
        ON CONFLICT (id) DO NOTHING
    `)
    if err != nil {
        log.Fatal("failed to insert test wallet:", err)
    }

    fmt.Println("Wallets table ready and test wallet inserted")
}


func (s *Storage) Get(id string) (*Wallet, error) {
	var w Wallet

	err := s.db.QueryRow(
		`SELECT id, balance FROM wallets WHERE id = $1`,
		id,
	).Scan(&w.ID, &w.Amount)

	if err == sql.ErrNoRows {
		return nil, errors.New("wallet not found")
	}

	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}

	return &w, nil
}


func (s *Storage) Deposit(id string, amount int64) (*Wallet, error) {
	var w Wallet

	err := s.db.QueryRow(
		`UPDATE wallets
		 SET balance = balance + $1
		 WHERE id = $2
		 RETURNING id, balance`,
		amount, id,
	).Scan(&w.ID, &w.Amount)

	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (s *Storage) Withdraw(id string, amount int64) (*Wallet, error) {
	var w Wallet

	err := s.db.QueryRow(
		`UPDATE wallets
		 SET balance = balance - $1
		 WHERE id = $2 AND balance >= $1
		 RETURNING id, balance`,
		amount, id,
	).Scan(&w.ID, &w.Amount)

	if err == sql.ErrNoRows {
		return nil, errors.New("not enough money")
	}

	if err != nil {
		return nil, err
	}

	return &w, nil
}
