package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

const (
	LogInfoStartBatchDb = "INFO.START.SetupPostgres\n"
	LogInfoEndBatchDb   = "INFO.END.SetupPostgres\n"
	LogErrorBatchDb     = "FATAL.SetupPostgres | reason: %s | error: %w\n"
)

func SetupPostgres() *sql.DB {
	log.Printf(LogInfoStartBatchDb)

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatalf(LogErrorBatchDb, "failed to open database", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf(LogErrorBatchDb, "failed to connect to database", err)
	}

	log.Printf(LogInfoEndBatchDb)
	return db
}
