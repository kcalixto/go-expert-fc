package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
)

func main() {
	http.HandleFunc("/temperature", getTemperature)
	http.ListenAndServe(":8080", nil)
}

func getTemperature(w http.ResponseWriter, r *http.Request) {
	// review warning
	if os.Getenv("WEATHER_API_KEY") == "" {
		w.Write([]byte("Conforme solicitado, a aplicação possui um endpoint no GCP para ser executada e não deve rodar localmente, uma vez que a API de clima requer uma chave de acesso que não deve ser exposta. Favor acessar o endpoint descrito no readme.md: https://cloudrun-goexpert-843349195325.southamerica-east1.run.app/temperature?cep=01001000"))
		return
	}

	cep := r.URL.Query().Get("cep")
	req := &Request{CEP: cep}
	if err := req.Validate(); err != nil {
		err := fmt.Errorf("invalid zipcode. To use this endpoint, you need to provide a valid zipcode, example: {url}/temperature?cep=01001000")
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	addr, err := callViaCep(cep)
	if err != nil {
		err = fmt.Errorf("cannot find zipcode: %w", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if addr.Localidade == "" {
		err = fmt.Errorf("cannot find zipcode")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	weatherApiResponse, err := callWeatherApi(addr)
	if err != nil {
		err = fmt.Errorf("cannot find weather: %w", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res := &Response{
		Celsius:    weatherApiResponse.Current.TempC,
		Fahrenheit: weatherApiResponse.Current.TempF,
		Kelvin:     weatherApiResponse.Current.TempC + 273.15,
	}

	responseBodyBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "cannot marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responseBodyBytes); err != nil {
		http.Error(w, "cannot write response", http.StatusInternalServerError)
		return
	}
}

type Request struct {
	CEP string `json:"cep" validate:"required,min=8,max=8"`
}

func (r *Request) Validate() error {
	return validator.New().Struct(r)
}

type Response struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

type WeatherApiResponse struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch float64 `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch float64 `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            float64 `json:"is_day"`
		Condition        struct {
			Text string  `json:"text"`
			Icon string  `json:"icon"`
			Code float64 `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree float64 `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   float64 `json:"humidity"`
		Cloud      float64 `json:"cloud"`
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
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

func callWeatherApi(addr *ViaCEPApiResponse) (*WeatherApiResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.weatherapi.com/v1/current.json", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("q", addr.Localidade)
	query.Set("key", os.Getenv("WEATHER_API_KEY"))
	req.URL.RawQuery = query.Encode()

	println(req.URL.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var weatherApiResponse WeatherApiResponse
	err = json.Unmarshal(responseBodyBytes, &weatherApiResponse)
	if err != nil {
		return nil, err
	}

	return &weatherApiResponse, nil
}

type ViaCEPApiResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func callViaCep(cep string) (*ViaCEPApiResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://viacep.com.br/ws/"+cep+"/json", nil)
	if err != nil {
		return nil, err
	}

	println(req.URL.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var viaCEPApiResponse ViaCEPApiResponse
	err = json.Unmarshal(responseBodyBytes, &viaCEPApiResponse)
	if err != nil {
		return nil, err
	}

	return &viaCEPApiResponse, nil
}
