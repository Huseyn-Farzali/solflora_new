package main

import (
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
	insertQueue chan<- DbEntry,
	notifyQueue chan<- DbEntry,
) {
	log.Printf("INFO.ESP_POLLING.START\n")
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			spMap := spState.GetAll()

			log.Printf(LogInfoStartBatchPolling, spMap)

			signalMap, err := FetchFromESP(spMap)
			if err != nil {
				log.Printf(LogErrorBatchPolling, spMap, "fetchFromESP failed", err)
				continue
			}

			for _, signal := range signalMap {
				select {
				case insertQueue <- signal:
				default:
					log.Println("⚠️ insertQueue full — dropping signal:", signal)
				}

				select {
				case notifyQueue <- signal:
				default:
					log.Println("⚠️ notifyQueue full — dropping signal:", signal)
				}
			}

			log.Printf(LogInfoEndBatchPolling, spMap)
		}
	}()

	log.Printf("INFO.ESP_POLLING.END\n")
}
