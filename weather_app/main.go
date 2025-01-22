package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Forecast struct {
	Temperature float32 `json:"temperature"`
	Unit        string
}

type WeatherProvider interface {
	GetForecast(location string) (*Forecast, error)
}

type VoyagerSpaceProbe struct {
	temps map[string]Forecast
}

func NewVoyagerSpaceProbe(temps map[string]Forecast) VoyagerSpaceProbe {
	return VoyagerSpaceProbe{temps}
}

func (v VoyagerSpaceProbe) GetForecast(location string) (*Forecast, error) {

	f, exists := v.temps[location]

	if !exists {
		return nil, fmt.Errorf("unknown location")
	}

	return &f, nil
}

func ServeWeather(wp WeatherProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		param := r.URL.Path

		// Check if a parameter exists
		if param == "" {
			http.Error(w, "Parameter not provided", http.StatusBadRequest)
			return
		}

		f, err := wp.GetForecast(param)

		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("error getting forecast for %s: %s", param, err),
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// write response back without allocating bytes
		if err := json.NewEncoder(w).Encode(f); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusInternalServerError)
			return
		}
	}
}

func main() {

	v := NewVoyagerSpaceProbe(map[string]Forecast{
		"CMB": {
			Temperature: -270.45,
			Unit:        "°C",
		},
	})

	omp := OpenMeteoProvider{}

	handlers := []struct {
		Path     string
		Provider WeatherProvider
	}{
		{
			Path:     "/forecast_voyager/",
			Provider: v,
		},
		{
			Path:     "/forecast_earth/",
			Provider: omp,
		},
	}

	for _, h := range handlers {

		log.Printf("registering path: %s", h.Path)
		path := h.Path
		handler := http.StripPrefix(path, ServeWeather(h.Provider))

		http.Handle(path, handler)

	}

	log.Println("listenning on port 8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatalf("error with server: %s", err)
	}
}
