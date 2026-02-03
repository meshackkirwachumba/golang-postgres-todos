package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectToDatabase(databaseURL string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Println("Error parsing database URL:", err)
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Println("Error creating connection pool:", err)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		log.Println("Error pinging database:", err)
		pool.Close()
		return nil, err
	}

	log.Println("Successfully connected to the database")
	
	return pool, nil
}