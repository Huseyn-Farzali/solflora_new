package main

import (
	"database/sql"
	"log"
	"time"
)

const (
	LogInfoStartBatchPolling = "INFO.START.StartESPPolling with spMap: %+v\n"
	LogInfoEndBatchPolling   = "INFO.END.StartESPPolling with spMap: %+v\n"
	LogErrorBatchPolling     = "ERROR.StartESPPolling with spMap: %+v | reason: %s | error: %w\n"
)

func StartESPPolling(
	interval time.Duration,
	spState *SPState,
	tuneState *TuneState,
	db *sql.DB,
) {
	log.Printf("INFO.START.StartEspPollingConfiguration\n")
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			spMap := spState.GetAll()
			tuneMap := tuneState.GetAll()

			log.Printf(LogInfoStartBatchPolling, spMap)

			dbEntries, err := FetchFromESPAndMapDbEntry(spMap, tuneMap)
			if err != nil {
				log.Printf(LogErrorBatchPolling, spMap, "fetchFromESPAndMapDbEntry failed", err)
				continue
			}

			err = insertBatchToDb(db, dbEntries)
			if err != nil {
				log.Printf(LogErrorBatchPolling, spMap, "insertBatchToDb failed", err)
				continue
			}

			log.Printf(LogInfoEndBatchPolling, spMap)
		}
	}()

	log.Printf("INFO.END.StartEspPollingConfiguration\n")
}
