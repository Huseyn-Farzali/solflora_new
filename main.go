package main

import (
	"log"
	"net/http"
)

func main() {
	dbConn := SetupPostgres()
	defer dbConn.Close()

	spState := NewSPState()
	tuneState := NewTuneState()

	StartESPPolling(EspPollingInterval, spState, tuneState, dbConn)

	http.HandleFunc("/api/data", HandleChartData(dbConn))
	http.HandleFunc("/api/setpoints", HandleSetPointUpdate(spState))
	http.HandleFunc("/api/update-pid", HandleTuningUpdate(tuneState))

	log.Println("🌱 Solflora backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
