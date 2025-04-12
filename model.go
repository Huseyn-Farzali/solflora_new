package main

import "time"

type PhysicalVariable string
type TuneVariable string

const (
	Temperature PhysicalVariable = "temperature"
	Humidity    PhysicalVariable = "humidity"
	Moisture    PhysicalVariable = "moisture"
)

const (
	KP TuneVariable = "kp"
	KI TuneVariable = "ki"
	KD TuneVariable = "kd"
)

type TuneProfile struct {
	KP float64 `json:"kp"`
	KI float64 `json:"ki"`
	KD float64 `json:"kd"`
}

type DbEntry struct {
	Variable  PhysicalVariable
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
