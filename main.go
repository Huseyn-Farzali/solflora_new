package main

import (
	"log"
	"net/http"
)

func main() {
	dbConn := SetupPostgres()
	defer dbConn.Close()

	spState := NewSPState()

	StartESPPolling(EspPollingInterval, spState, dbConn)

	// 6. HTTP routes
	http.HandleFunc("/api/sp", HandleChartData(dbConn))

	log.Println("🌱 Solflora backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
