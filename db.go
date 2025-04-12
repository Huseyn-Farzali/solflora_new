package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
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

func insertBatchToDb(db *sql.DB, entries []DbEntry) error {
	log.Printf("INFO.START.insertBatchToDb with entries: %v\n", entries)
	if len(entries) == 0 {
		return nil
	}

	const baseQuery = "INSERT INTO entries (variable, timestamp, sp, pv, co) VALUES "

	valueStrings := make([]string, 0, len(entries))
	valueArgs := make([]interface{}, 0, len(entries)*5)

	for i, entry := range entries {
		startIdx := i * 5
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
			startIdx+1, startIdx+2, startIdx+3, startIdx+4, startIdx+5))

		valueArgs = append(valueArgs,
			entry.Variable,
			entry.Timestamp,
			entry.SP,
			entry.PV,
			entry.CO,
		)
	}

	query := baseQuery + strings.Join(valueStrings, ",")
	_, err := db.Exec(query, valueArgs...)
	if err != nil {
		log.Printf("ERROR.insertBatchToDb with entries: %v\n", entries)
		return err
	}

	log.Printf("INFO.END.insertBatchToDb with entries: %v\n", entries)
	return nil
}
