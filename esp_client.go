package main

import (
	"bytes"
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

func FetchFromESP(spMap map[Variable]float64) (map[Variable]DbEntry, error) {
	log.Printf(LogInfoStartBatch, spMap)

	reqBody := make(ESPRequestBody)
	for variable, sp := range spMap {
		reqBody[variable] = struct {
			SP float64 `json:"SP"`
		}{SP: sp}
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf(LogErrorBatch, spMap, "json marshalling failed", err)
	}

	resp, err := http.Post(EspBaseUrl, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf(LogErrorBatch, spMap, "[POST] request to ESP failed", err)
	}
	defer resp.Body.Close()

	var response ESPResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf(LogErrorBatch, spMap, "failed to decode ESP response", err)
	}

	result := make(map[Variable]DbEntry)
	now := time.Now()

	for variable, data := range response {
		result[variable] = DbEntry{
			Variable:  variable,
			TimeStamp: now,
			SP:        spMap[variable],
			PV:        data.PV,
			CO:        data.CO,
		}
	}

	log.Printf(LogInfoEndBatch, spMap)
	return result, nil
}
