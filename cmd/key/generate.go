package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	rawKey, err := generateAPIKey()
	if err != nil {
		log.Fatal("failed to generate api key:", err)
	}

	hash := hashKey(rawKey)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	_, err = pool.Exec(
		ctx,
		`INSERT INTO api_keys (key_hash) VALUES ($1)`,
		hash,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("API key created successfully")
	fmt.Println()
	fmt.Println("RAW API KEY (store this safely):")
	fmt.Println(rawKey)
	fmt.Println()
	fmt.Println("This key will not be shown again.")
}

func generateAPIKey() (string, error) {
	b := make([]byte, 32) // 256-bit key
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return "sw_" + hex.EncodeToString(b), nil
}

func hashKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}
