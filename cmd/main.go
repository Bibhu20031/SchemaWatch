package main

import (
	"log"

	database "github.com/Bibhu20031/SchemaWatch/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	_, err := database.ConnectDB()
	if err != nil {
		log.Fatal("failed to connect database")
	}

	log.Println("service started")
}
