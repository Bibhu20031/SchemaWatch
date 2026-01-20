package main

import (
	"log"

	"github.com/Bibhu20031/SchemaWatch/internal/config"
)

func main() {
	cfg := config.Load()

	log.Printf("starting schemawatch | env=%s | port=%s",
		cfg.AppEnv,
		cfg.AppPort,
	)

}
