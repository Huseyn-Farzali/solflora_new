package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	LogInfoStartBatch = "INFO.START.fetchFromESP with SP map: %+v\n"
	LogInfoEndBatch   = "INFO.END.fetchFromESP with SP map: %+v\n"
	LogErrorBatch     = "ERROR.fetchFromESP with SP map: %+v | reason: %s\n | error: %v\n"
)

type ResponseData struct {
	PV float64 `json:"PV"`
	CO float64 `json:"CO"`
}

type ESPRequestBody map[PhysicalVariable]struct {
	SP float64 `json:"sp"`
	KP float64 `json:"kp"`
	KI float64 `json:"ki"`
	KD float64 `json:"kd"`
}

type ESPResponseBody map[PhysicalVariable]ResponseData

func FetchFromESPAndMapDbEntry(spMap map[PhysicalVariable]float64, tuneMap map[PhysicalVariable]TuneProfile) ([]DbEntry, error) {
	log.Printf(LogInfoStartBatch, spMap)

	reqBody := make(ESPRequestBody)
	for variable, _ := range spMap {
		reqBody[variable] = struct {
			SP float64 `json:"sp"`
			KP float64 `json:"kp"`
			KI float64 `json:"ki"`
			KD float64 `json:"kd"`
		}{SP: spMap[variable], KP: tuneMap[variable].KP, KI: tuneMap[variable].KI, KD: tuneMap[variable].KD}
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf(LogErrorBatch, spMap, "json marshalling failed", err)
		return nil, fmt.Errorf("json marshalling failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), EspPollingCallTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, EspBaseUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf(LogErrorBatch, spMap, "creating HTTP request failed", err)
		return nil, fmt.Errorf("creating HTTP request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf(LogErrorBatch, spMap, "[POST] request to ESP failed", err)
		return nil, fmt.Errorf("[POST] request to ESP failed: %w", err)
	}
	defer resp.Body.Close()

	var response ESPResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf(LogErrorBatch, spMap, "failed to decode ESP response", err)
		return nil, fmt.Errorf("failed to decode ESP response: %w", err)
	}

	result := make([]DbEntry, 3)
	now := time.Now()

	index := 0
	for variable, data := range response {
		result[index] = DbEntry{
			Variable:  variable,
			Timestamp: now,
			SP:        spMap[variable],
			PV:        data.PV,
			CO:        data.CO,
		}
		index++
	}

	log.Printf(LogInfoEndBatch, spMap)
	return result, nil
}
