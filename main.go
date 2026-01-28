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
	db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err != nil {
		log.Fatal(err)
	}

	storage := NewStorage(db)
	storage.InitDB()

	handler := NewHandler(storage)

	router := gin.Default()
	handler.RegisterRoutes(router)

	router.Run(":8080")
}
