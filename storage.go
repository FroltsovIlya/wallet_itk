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

type StorageInterface interface { //create interface for tests
	Get(id string) (*Wallet, error)
	Deposit(id string, amount int64) (*Wallet, error)
	Withdraw(id string, amount int64) (*Wallet, error)
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) InitDB() {
    connStr := fmt.Sprintf( //create string to connect to db
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

    var err error
    s.db, err = sql.Open("postgres", connStr) //open connection to db
    if err != nil {
        panic(err)
    }

    err = s.db.Ping() //dergaem db for get information is it ready
    if err != nil { //if db is not answer, waiting 5 seconds, may be it not started yet
        log.Println("waiting for database...")
        time.Sleep(5 * time.Second)
        }

    fmt.Println("Successfully connected to PostgreSQL")
        //here we create table in db if it not exists
    _, err = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS wallets (
            id TEXT PRIMARY KEY,
            balance BIGINT NOT NULL DEFAULT 0
        )
    `)
    if err != nil {
        log.Fatal("failed to create wallets table:", err)
    }
    //happy to duplicate code of creating objects in db
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
    //if no errors, table ready to work
    fmt.Println("Wallets table ready and test wallet inserted")
}

//get func for walllet
func (s *Storage) Get(id string) (*Wallet, error) {
	var w Wallet //create wallet

	err := s.db.QueryRow( //getting wallet from db by id
		`SELECT id, balance FROM wallets WHERE id = $1`,
		id,
	).Scan(&w.ID, &w.Amount)//getting data from returned wallet

	if err == sql.ErrNoRows { //compare error with no row to log
		return nil, errors.New("wallet not found")
	}

	if err != nil { //if error exist, show it
		return nil, fmt.Errorf("db error: %w", err)
	}

	return &w, nil //return our getted wallet
}


func (s *Storage) Deposit(id string, amount int64) (*Wallet, error) {
	var w Wallet

	err := s.db.QueryRow( //update wallet and get new
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

	err := s.db.QueryRow( //also here, like Deposit(). But we looking for amount of operation was less than summary on wallet
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
