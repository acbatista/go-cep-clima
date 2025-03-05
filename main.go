package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"

	"go-cep-clima/service"
)

type WeatherResponse struct {
	TempC     float64 `json:"temp_C"`
	TempF     float64 `json:"temp_F"`
	TempK     float64 `json:"temp_K"`
	Condition struct {
		Text string `json:"text"`
		Icon string `json:"icon"`
	} `json:"condition"`
	Location struct {
		City    string `json:"city"`
		State   string `json:"state"`
		Country string `json:"country"`
	} `json:"location"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func isValidCEP(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}

func handleWeather(weatherSvc *service.WeatherService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		cep := r.URL.Query().Get("cep")
		if !isValidCEP(cep) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid zipcode"})
			return
		}

		location, err := weatherSvc.GetLocationByCEP(cep)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "can not find zipcode"})
			return
		}

		weatherData, err := weatherSvc.GetWeatherByCity(location.Localidade, location.UF, location.CEP)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "error fetching weather data"})
			return
		}

		tempC := weatherData.Current.TempC
		response := WeatherResponse{
			TempC: tempC,
			TempF: service.CelsiusToFahrenheit(tempC),
			TempK: service.CelsiusToKelvin(tempC),
			Condition: struct {
				Text string `json:"text"`
				Icon string `json:"icon"`
			}{
				Text: weatherData.Current.Condition.Text,
				Icon: weatherData.Current.Condition.Icon,
			},
			Location: struct {
				City    string `json:"city"`
				State   string `json:"state"`
				Country string `json:"country"`
			}{
				City:    weatherData.Location.Name,
				State:   weatherData.Location.Region,
				Country: weatherData.Location.Country,
			},
		}

		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	weatherSvc := service.NewWeatherService(weatherAPIKey)
	http.HandleFunc("/weather", handleWeather(weatherSvc))

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
