package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	r := NewRequest()
	s := NewStorage()
	s.FillStorage()

	router.POST("/api/v1/wallet", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
		
		c.BindJSON(&r)
		id := r.ID
		amount := r.Amount

		out, err := s.Get(id)
		if err == nil {
			err = r.Transaction(amount, &out)
			if err == nil {
				c.JSON(http.StatusOK, out)
			} else {
				c.JSON(http.StatusOK, gin.H{
					"badtransaction": "baaad",
				})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"badrequest": "baaad",
			})
		}
	})

	router.GET("/api/v1/wallets/:id", func(c *gin.Context) {
		id := c.Param("id")
		out, err := s.Get(id)
		if err == nil {
			c.JSON(http.StatusOK, out)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"badrequest": "baaad",
			})
		}
	})

	router.Run(":8080")
}
