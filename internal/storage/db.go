package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func ConnectDB() (*pgxpool.Pool, error) {
	if pool != nil {
		return pool, nil
	}

	dsn := os.Getenv("DATABASE_URL")

	p, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Printf("unable to connect to database: %v\n", err)
		return nil, err
	}

	if err := p.Ping(context.Background()); err != nil {
		log.Printf("unable to ping database: %v\n", err)
		return nil, err
	}

	pool = p
	log.Println("connected to database")
	return pool, nil
}
