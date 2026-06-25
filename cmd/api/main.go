package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"takehome/internal/db"
	"takehome/internal/models"
)

func main() {
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		rows, err := conn.Query("SELECT sku, name, current_stock FROM products ORDER BY sku")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var p models.Product
			if err := rows.Scan(&p.SKU, &p.Name, &p.CurrentStock); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})

	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse SKU from path: /api/products/{sku}/movements
		path := r.URL.Path
		const prefix = "/api/products/"
		if len(path) <= len(prefix) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		rest := path[len(prefix):]
		// rest should be "{sku}/movements"
		var sku string
		for i, c := range rest {
			if c == '/' {
				sku = rest[:i]
				rest = rest[i:]
				break
			}
		}

		if sku == "" || rest != "/movements" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		pageStr := r.URL.Query().Get("page")
		page := 1
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
		
		limitStr := r.URL.Query().Get("limit")
		limit := 100
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
		
		offset := (page - 1) * limit

		rows, err := conn.Query("SELECT event_id, sku, type, quantity, occurred_at FROM movements WHERE sku = $1 ORDER BY occurred_at DESC LIMIT $2 OFFSET $3", sku, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var events []models.Event
		for rows.Next() {
			var e models.Event
			if err := rows.Scan(&e.EventID, &e.SKU, &e.Type, &e.Quantity, &e.OccurredAt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			events = append(events, e)
		}
		if events == nil {
			events = []models.Event{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	})

	log.Println("API listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
