package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-cep-clima/service"
)

func TestIsValidCEP(t *testing.T) {
	tests := []struct {
		name string
		cep  string
		want bool
	}{
		{"valid CEP", "12345678", true},
		{"invalid CEP with letters", "1234567a", false},
		{"invalid CEP too short", "1234567", false},
		{"invalid CEP too long", "123456789", false},
		{"invalid CEP with special chars", "12345-678", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidCEP(tt.cep); got != tt.want {
				t.Errorf("isValidCEP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandleWeather_InvalidCEP(t *testing.T) {
	weatherSvc := service.NewWeatherService("dummy-key")
	handler := handleWeather(weatherSvc)

	req := httptest.NewRequest("GET", "/weather?cep=invalid", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}

	var response ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Message != "invalid zipcode" {
		t.Errorf("Expected message 'invalid zipcode', got '%s'", response.Message)
	}
}

func TestWeatherResponse_Structure(t *testing.T) {
	response := WeatherResponse{
		TempC: 25.0,
		TempF: 77.0,
		TempK: 298.15,
		Condition: struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
		}{
			Text: "Sunny",
			Icon: "//cdn.weatherapi.com/weather/64x64/day/113.png",
		},
		Location: struct {
			City    string `json:"city"`
			State   string `json:"state"`
			Country string `json:"country"`
		}{
			City:    "São Paulo",
			State:   "São Paulo",
			Country: "Brazil",
		},
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var decoded WeatherResponse
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if decoded.TempC != 25.0 || decoded.TempF != 77.0 || decoded.TempK != 298.15 {
		t.Errorf("Temperature values don't match expected values")
	}

	if decoded.Condition.Text != "Sunny" || decoded.Condition.Icon == "" {
		t.Errorf("Condition values don't match expected values")
	}

	if decoded.Location.City != "São Paulo" || decoded.Location.State != "São Paulo" || decoded.Location.Country != "Brazil" {
		t.Errorf("Location values don't match expected values")
	}
}
