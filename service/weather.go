package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ViaCEPResponse struct {
	CEP        string `json:"cep"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	Error      bool   `json:"erro"`
}

type WeatherAPIResponse struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int64   `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		TempF     float64 `json:"temp_f"`
		IsDay     int     `json:"is_day"`
		Condition struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		WindchillC float64 `json:"windchill_c"`
		WindchillF float64 `json:"windchill_f"`
		HeatindexC float64 `json:"heatindex_c"`
		HeatindexF float64 `json:"heatindex_f"`
		DewpointC  float64 `json:"dewpoint_c"`
		DewpointF  float64 `json:"dewpoint_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		UV         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

type WeatherService struct {
	weatherAPIKey string
	baseURL       string
}

func NewWeatherService(weatherAPIKey string) *WeatherService {
	return &WeatherService{
		weatherAPIKey: weatherAPIKey,
		baseURL:       "http://api.weatherapi.com/v1",
	}
}

func (s *WeatherService) GetLocationByCEP(cep string) (*ViaCEPResponse, error) {
	resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("can not find zipcode")
	}

	var viaCEP ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCEP); err != nil {
		return nil, err
	}

	if viaCEP.Error {
		return nil, fmt.Errorf("can not find zipcode")
	}

	return &viaCEP, nil
}

func (s *WeatherService) GetWeatherByCity(city, state string, cep string) (*WeatherAPIResponse, error) {

	url := fmt.Sprintf("%s/current.json?key=%s&q=%s",
		s.baseURL, s.weatherAPIKey, cep)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get weather data: status %d", resp.StatusCode)
	}

	var weather WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("failed to decode weather data: %v", err)
	}

	return &weather, nil
}

func CelsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func CelsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}
