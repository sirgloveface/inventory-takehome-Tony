package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"takehome/internal/db"
	"takehome/internal/models"
)

func main() {
	// 1. Context with cancellation for SIGINT/SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Received termination signal, shutting down gracefully...")
		cancel()
	}()

	// 2. Connect to Database
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	// 3. Initialize Schema
	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		log.Fatalf("Failed to read schema.sql: %v", err)
	}
	if _, err := conn.ExecContext(ctx, string(schema)); err != nil {
		log.Fatalf("Failed to execute schema.sql: %v", err)
	}
	log.Println("Schema initialized successfully.")

	// 4. Load Products
	if err := loadProducts(ctx, conn, "data/products.csv"); err != nil {
		log.Fatalf("Failed to load products: %v", err)
	}
	log.Println("Products loaded successfully.")

	// 5. Ingest Events
	files, err := filepath.Glob("data/events/part-*.ndjson")
	if err != nil {
		log.Fatalf("Failed to find event files: %v", err)
	}

	var wg sync.WaitGroup
	// Create a channel to distribute files to workers
	fileCh := make(chan string, len(files))
	for _, f := range files {
		fileCh <- f
	}
	close(fileCh)

	// Since max conns is 10, we'll use 5 workers so we don't starve connections.
	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, &wg, fileCh, conn)
	}

	wg.Wait()
	log.Println("Ingestion process completed.")
}

func loadProducts(ctx context.Context, db *sql.DB, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for i, row := range records {
		if i == 0 { // Skip header
			continue
		}
		if len(row) < 2 {
			continue
		}
		sku, name := row[0], row[1]
		_, err := db.ExecContext(ctx, `
			INSERT INTO products (sku, name, current_stock)
			VALUES ($1, $2, 0)
			ON CONFLICT (sku) DO NOTHING
		`, sku, name)
		if err != nil {
			return fmt.Errorf("failed to insert product %s: %w", sku, err)
		}
	}
	return nil
}

func worker(ctx context.Context, wg *sync.WaitGroup, fileCh <-chan string, db *sql.DB) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case file, ok := <-fileCh:
			if !ok {
				return // No more files
			}
			processFile(ctx, db, file)
		}
	}
}

func processFile(ctx context.Context, db *sql.DB, file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Printf("Error opening file %s: %v", file, err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		// Check for context cancellation before processing next line
		if err := ctx.Err(); err != nil {
			log.Printf("Context cancelled, stopping file %s at line %d", file, lineNum)
			return
		}

		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var event models.Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			log.Printf("Invalid JSON in file %s:%d - %v", file, lineNum, err)
			continue
		}

		// Validation
		if event.Quantity <= 0 {
			log.Printf("Invalid quantity in file %s:%d - %d", file, lineNum, event.Quantity)
			continue
		}
		if event.Type != "IN" && event.Type != "OUT" {
			log.Printf("Invalid type in file %s:%d - %s", file, lineNum, event.Type)
			continue
		}

		// Insert with idempotency and update stock atomically
		query := `
		WITH inserted AS (
			INSERT INTO movements (event_id, sku, type, quantity, occurred_at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (event_id) DO NOTHING
			RETURNING sku, type, quantity
		)
		UPDATE products p
		SET current_stock = current_stock + CASE WHEN i.type = 'IN' THEN i.quantity ELSE -i.quantity END
		FROM inserted i
		WHERE p.sku = i.sku;
		`
		
		_, err = db.ExecContext(ctx, query, event.EventID, event.SKU, event.Type, event.Quantity, event.OccurredAt)
		if err != nil {
			// It might fail if SKU doesn't exist due to FK constraint
			log.Printf("Failed to insert event %s in file %s:%d - %v", event.EventID, file, lineNum, err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file %s: %v", file, err)
	} else {
		log.Printf("Finished processing %s", file)
	}
}
