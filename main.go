package main

import (
	"log"
	"net/http"
)

func main() {
	dbConn := SetupPostgres()
	defer dbConn.Close()

	spState := NewSPState()

	insertQueue := make(chan DbEntry, 50)
	notifyQueue := make(chan DbEntry, 50)

	for i := 0; i < 4; i++ {
		go func() {
			for signal := range insertQueue {
				InsertToDB(dbConn, signal)
			}
		}()
	}

	go func() {
		for signal := range notifyQueue {
			Broadcast(signal) // to WebSocket clients
		}
	}()

	StartESPPolling(PollingInterval, spState, insertQueue, notifyQueue)

	// 6. HTTP routes
	http.HandleFunc("/api/sp", HandleSPUpdate(spState)) // accepts {"temperature": {"SP": 45.3}} etc.
	//http.HandleFunc("/ws", HandleWebSocket())           // optional websocket endpoint

	log.Println("🌱 Solflora backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
