package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// OpenMeteoProvider implements the WeatherProvider interface using the Open-Meteo API.
type OpenMeteoProvider struct{}

// GetForecast fetches the weather forecast for a given location.
func (omp OpenMeteoProvider) GetForecast(location string) (*Forecast, error) {
	// Mock coordinates for simplicity (use geocoding APIs for dynamic location resolution)
	coordinates := map[string][2]float64{
		"london":  {51.5074, -0.1278},
		"newyork": {40.7128, -74.0060},
		"paris":   {48.8566, 2.3522},
	}

	// Lookup coordinates for the location
	coords, ok := coordinates[location]
	if !ok {
		return nil, fmt.Errorf("location not found: %s", location)
	}

	// Open-Meteo API endpoint
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current_weather=true",
		coords[0], coords[1],
	)

	// Send GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching weather data: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 response codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Parse the JSON response
	var data struct {
		CurrentWeather struct {
			Temperature float64 `json:"temperature"`
		} `json:"current_weather"`
		CurrentWeatherUnits struct {
			Temperature string `json:"temperature"`
		} `json:"current_weather_units"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Create a Forecast object from the API data
	forecast := &Forecast{
		Temperature: float32(data.CurrentWeather.Temperature),
		Unit:        data.CurrentWeatherUnits.Temperature,
	}
	return forecast, nil
}
