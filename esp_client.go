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

type ESPRequestBody map[Variable]struct {
	SP float64 `json:"SP"`
}

type ESPResponseBody map[Variable]ResponseData

func FetchFromESPAndMapDbEntry(spMap map[Variable]float64) ([]DbEntry, error) {
	log.Printf(LogInfoStartBatch, spMap)

	reqBody := make(ESPRequestBody)
	for variable, sp := range spMap {
		reqBody[variable] = struct {
			SP float64 `json:"SP"`
		}{SP: sp}
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
