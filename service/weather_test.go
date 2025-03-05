package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name    string
		celsius float64
		want    float64
	}{
		{"zero celsius", 0, 32},
		{"positive celsius", 25, 77},
		{"negative celsius", -10, 14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CelsiusToFahrenheit(tt.celsius)
			if got != tt.want {
				t.Errorf("CelsiusToFahrenheit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		name    string
		celsius float64
		want    float64
	}{
		{"zero celsius", 0, 273.15},
		{"positive celsius", 25, 298.15},
		{"negative celsius", -10, 263.15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CelsiusToKelvin(tt.celsius)
			if got != tt.want {
				t.Errorf("CelsiusToKelvin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeatherService_GetWeatherByCity(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é uma requisição GET
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Verificar se os parâmetros necessários estão presentes
		q := r.URL.Query()
		if q.Get("key") == "" || q.Get("q") == "" {
			t.Errorf("Missing required parameters")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Resposta simplificada mas válida
		response := `{
			"location": {
				"name": "São Paulo",
				"region": "São Paulo",
				"country": "Brazil"
			},
			"current": {
				"temp_c": 25.0,
				"temp_f": 77.0,
				"condition": {
					"text": "Sunny",
					"icon": "//cdn.weatherapi.com/weather/64x64/day/113.png",
					"code": 1000
				}
			}
		}`

		w.Write([]byte(response))
	}))
	defer server.Close()

	// Criar serviço usando o servidor mock
	svc := &WeatherService{
		weatherAPIKey: "dummy-key",
		baseURL:       server.URL, // Usar URL do servidor mock
	}

	// Executar o teste
	weather, err := svc.GetWeatherByCity("São Paulo", "SP", "01001000")
	if err != nil {
		t.Fatalf("GetWeatherByCity() error = %v", err)
	}

	// Verificações
	if weather.Current.TempC != 25.0 {
		t.Errorf("Expected temperature 25.0°C, got %v°C", weather.Current.TempC)
	}

	if weather.Location.Name != "São Paulo" {
		t.Errorf("Expected city São Paulo, got %v", weather.Location.Name)
	}

	if weather.Current.Condition.Text != "Sunny" {
		t.Errorf("Expected condition Sunny, got %v", weather.Current.Condition.Text)
	}
}

func TestWeatherService_GetLocationByCEP(t *testing.T) {
	// Mock server para simular a ViaCEP API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ViaCEPResponse{
			CEP:        "01001000",
			Localidade: "São Paulo",
			UF:         "SP",
			Error:      false,
		}

		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Criar serviço com o servidor mock
	svc := NewWeatherService("dummy-key")

	// Testar a função
	location, err := svc.GetLocationByCEP("01001000")
	if err != nil {
		t.Fatalf("GetLocationByCEP() error = %v", err)
	}

	if location.Localidade != "São Paulo" {
		t.Errorf("Expected city São Paulo, got %v", location.Localidade)
	}

	if location.UF != "SP" {
		t.Errorf("Expected state SP, got %v", location.UF)
	}
}
