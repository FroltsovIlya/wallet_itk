package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
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
