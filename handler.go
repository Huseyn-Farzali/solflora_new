package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type SimplifiedEntry struct {
	SP float64 `json:"sp"`
	PV float64 `json:"pv"`
	CO float64 `json:"co"`
}

func HandleChartData(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("INFO.START.HandleChartData\n")

		intervalStr := r.URL.Query().Get("interval")
		if intervalStr == "" {
			http.Error(w, "missing 'interval' query parameter", http.StatusBadRequest)
			return
		}

		duration, err := time.ParseDuration(intervalStr)
		if err != nil {
			http.Error(w, "invalid interval format, use e.g. 10m or 1h", http.StatusBadRequest)
			return
		}

		variables := []PhysicalVariable{
			Temperature,
			Humidity,
			Moisture,
		}

		responseData := make(map[PhysicalVariable][]SimplifiedEntry)

		query := `
			SELECT sp, pv, co
			FROM signals
			WHERE variable = $1 AND timestamp >= NOW() - $2
			ORDER BY timestamp ASC;
		`

		for _, variable := range variables {
			rows, err := db.Query(query, variable, duration)
			if err != nil {
				log.Printf("ERROR.query.%s: %v\n", variable, err)
				http.Error(w, "failed to query "+string(variable), http.StatusInternalServerError)
				return
			}

			var entries []SimplifiedEntry
			for rows.Next() {
				var entry SimplifiedEntry
				if err := rows.Scan(&entry.SP, &entry.PV, &entry.CO); err != nil {
					log.Printf("ERROR.scan.%s: %v\n", variable, err)
					http.Error(w, "failed to scan "+string(variable), http.StatusInternalServerError)
					rows.Close()
					return
				}
				entries = append(entries, entry)
			}
			rows.Close()
			responseData[variable] = entries
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(responseData); err != nil {
			log.Printf("ERROR.encode.response: %v\n", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}

		log.Printf("INFO.END.HandleChartData\n")
	}
}

func HandleSetPointUpdate(spState *SPState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var updates map[PhysicalVariable]float64
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		for variable, sp := range updates {
			switch variable {
			case Temperature, Moisture:
				spState.Set(variable, sp)
			default: // Ignore humidity
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}

}

func HandleTuningUpdate(tuneState *TuneState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var updates map[PhysicalVariable]TuneProfile
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		for variable, profile := range updates {
			switch variable {
			case Temperature, Moisture:
				tuneState.Set(variable, profile)
			default: // Ignore humidity
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
