package main

import (
	"database/sql"
	"log"
	"os"
	"time"
	
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DB_DSN")) //open conn to Postgres
	db.SetMaxOpenConns(50) //configure connection
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err != nil { //if error in connection get log
		log.Fatal(err)
	}

	storage := NewStorage(db) //creating storage to handle
	storage.InitDB() //init db, creating table and adding basic wallets

	handler := NewHandler(storage) //creating handler

	router := gin.Default()
	handler.RegisterRoutes(router) //adding endpoints in handler

	router.Run(":8080")
}
