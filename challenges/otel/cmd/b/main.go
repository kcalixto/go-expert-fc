package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var tracer trace.Tracer

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("service-b"),
		),
	)
	if err != nil {
		fmt.Println("failed to create resource: %w", err)
		return
	}
	ctx, cancel = context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "otel-collector:4317",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		fmt.Println("failed to create gRPC connection to collector: %w", err)
		return
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		fmt.Println("failed to create trace exporter: %w", err)
		return
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	tracer = otel.Tracer("github.com/kcalixto/go-expert-fc/challanges/otel/cmd/b")

	// Create service
	ctx = context.Background()
	http.HandleFunc("/temperature", getTemperature)
	http.ListenAndServe(":8081", nil)

	fmt.Println("Service A is running...")
	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason...")
	}

	err = tracerProvider.Shutdown(ctx)
	if err != nil {
		log.Fatalf("failed to shutdown provider: %v", err)
		return
	}
}

func getTemperature(w http.ResponseWriter, r *http.Request) {
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	ctx, span := tracer.Start(ctx, "get-temperature-service-b")
	defer span.End()

	cep := r.URL.Query().Get("cep")
	req := &Request{CEP: cep}
	if err := req.Validate(); err != nil {
		err := fmt.Errorf("invalid zipcode. To use this endpoint, you need to provide a valid zipcode, example: {url}/temperature?cep=01001000")
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	addr, err := callViaCep(ctx, cep)
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

	weatherApiResponse, err := callWeatherApi(ctx, addr)
	if err != nil {
		err = fmt.Errorf("cannot find weather: %w", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res := &Response{
		DescricaoDesafio: `O desafio é composto por dois serviços, o serviço A e o serviço B, que integram com APIs externas e enviam traces para o OTEL Collector que está rodando em um container Docker.`,
		City:             weatherApiResponse.Location.Name,
		Celsius:          weatherApiResponse.Current.TempC,
		Fahrenheit:       weatherApiResponse.Current.TempF,
		Kelvin:           weatherApiResponse.Current.TempC + 273.15,
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
	DescricaoDesafio string  `json:"descricao_desafio"`
	City             string  `json:"city"`
	Celsius          float64 `json:"temp_C"`
	Fahrenheit       float64 `json:"temp_F"`
	Kelvin           float64 `json:"temp_K"`
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

func callWeatherApi(ctx context.Context, addr *ViaCEPApiResponse) (*WeatherApiResponse, error) {
	ctx, span := tracer.Start(ctx, "waether-api-request")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.weatherapi.com/v1/current.json", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("q", addr.Localidade)
	query.Set("key", "c25068cc19714f8182d200448241611")
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

func callViaCep(ctx context.Context, cep string) (*ViaCEPApiResponse, error) {
	ctx, span := tracer.Start(ctx, "viacep-api-request")
	defer span.End()

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
