package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"vertigo/pkg/broker"
	"vertigo/pkg/config"

	_ "modernc.org/sqlite"
)

func main() {
	fmt.Println("🚀 Starting Vertigo (CLI Demo)...")

	// 1. Load Configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Fatal: Failed to load config.yaml: %v", err)
	}

	// 2. Initialize Facade
	b, err := broker.NewBroker(cfg)
	if err != nil {
		log.Fatalf("Fatal: Failed to initialize Vertigo: %v", err)
	}

	// 2. Setup Schema for Demo
	_, err = b.DB.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, email TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// Insert dummy data
	_, _ = b.DB.Exec("INSERT INTO users (name, email) VALUES (?, ?), (?, ?)",
		"Alice", "alice@example.com", "Bob", "bob@example.com")

	// 3. Dispatch Query (SQL Abstraction)
	fmt.Println("📡 Dispatching SQL query: SELECT id FROM users...")
	ctx := context.Background()
	data, err := b.Dispatch(ctx, "SELECT id FROM users", "public_channel")
	if err != nil {
		log.Fatalf("Dispatch failed: %v", err)
	}

	fmt.Printf("✅ Dispatch complete! Captured Payload: %s\n", string(data))

	// Cleanup for demo
	os.Remove("demo.db")
}
