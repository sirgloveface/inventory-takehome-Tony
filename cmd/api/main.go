package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"takehome/internal/db"
	"takehome/internal/models"
)

func main() {
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.GET("/products", func(c *gin.Context) {
			rows, err := conn.Query("SELECT sku, name, current_stock FROM products ORDER BY sku")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer rows.Close()

			var products []models.Product
			for rows.Next() {
				var p models.Product
				if err := rows.Scan(&p.SKU, &p.Name, &p.CurrentStock); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				products = append(products, p)
			}
			if products == nil {
				products = []models.Product{}
			}

			c.JSON(http.StatusOK, products)
		})

		api.GET("/products/:sku/movements", func(c *gin.Context) {
			sku := c.Param("sku")

			pageStr := c.DefaultQuery("page", "1")
			page, err := strconv.Atoi(pageStr)
			if err != nil || page < 1 {
				page = 1
			}

			limitStr := c.DefaultQuery("limit", "100")
			limit, err := strconv.Atoi(limitStr)
			if err != nil || limit < 1 || limit > 1000 {
				limit = 100
			}

			offset := (page - 1) * limit

			rows, err := conn.Query("SELECT event_id, sku, type, quantity, occurred_at FROM movements WHERE sku = $1 ORDER BY occurred_at DESC LIMIT $2 OFFSET $3", sku, limit, offset)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer rows.Close()

			var events []models.Event
			for rows.Next() {
				var e models.Event
				if err := rows.Scan(&e.EventID, &e.SKU, &e.Type, &e.Quantity, &e.OccurredAt); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				events = append(events, e)
			}
			if events == nil {
				events = []models.Event{}
			}

			c.JSON(http.StatusOK, events)
		})
	}

	log.Println("API listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
