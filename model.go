package main

import "time"

type Variable string

const (
	Temperature Variable = "temperature"
	Humidity    Variable = "humidity"
	Moisture    Variable = "moisture"
)

type DbEntry struct {
	Variable  Variable
	Timestamp time.Time
	SP        float64
	PV        float64
	CO        float64
}

const (
	EspBaseUrl            string        = "http://esp32.local/api/"
	EspPollingInterval    time.Duration = 1 * time.Second
	EspPollingCallTimeout time.Duration = 3 * time.Second
)
