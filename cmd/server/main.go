package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tajri15/go-pulse-monitoring/internal/api"
	"github.com/tajri15/go-pulse-monitoring/internal/db"
)

func main() {
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "postgresql://user:password@localhost:5432/gopulse?sslmode=disable"
	}

	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}