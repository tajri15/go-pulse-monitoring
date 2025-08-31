package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tajri15/go-pulse-monitoring/internal/api"
	"github.com/tajri15/go-pulse-monitoring/internal/db"
	"github.com/tajri15/go-pulse-monitoring/internal/worker"
)

func main() {
	// Mengambil connection string dari environment variable, dengan fallback untuk local dev
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "postgresql://user:password@localhost:5432/gopulse?sslmode=disable"
	}

	// Membuat koneksi ke database PostgreSQL
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	// Membuat instance dari store database
	store := db.NewStore(conn)

	// Inisialisasi checker
	checker := worker.NewChecker(store)
	// Jalankan checker di background sebagai goroutine agar tidak memblokir server API
	go checker.Start()

	// Inisialisasi dan jalankan server API di foreground
	server := api.NewServer(store)
	err = server.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}