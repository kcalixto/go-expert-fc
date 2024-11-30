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
			semconv.ServiceName("service-a"),
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

	tracer = otel.Tracer("github.com/kcalixto/go-expert-fc/challanges/otel/cmd/a")

	// Create service
	ctx = context.Background()
	http.HandleFunc("/temperature", getTemperature)
	http.ListenAndServe(":8080", nil)

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
	ctx, span := tracer.Start(r.Context(), "get-temperature-service-a")
	defer span.End()

	type Schema struct {
		CEP string `json:"cep"`
	}

	var schema Schema
	if err := json.NewDecoder(r.Body).Decode(&schema); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "requisição inválida", http.StatusInternalServerError)
		return
	}

	if len(schema.CEP) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	temperatureResponse, err := callTemperatureService(ctx, schema.CEP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(temperatureResponse)
	w.WriteHeader(http.StatusOK)
}

func callTemperatureService(ctx context.Context, cep string) ([]byte, error) {
	ctx, span := tracer.Start(ctx, "call-temperature-service")
	defer span.End()

	url := fmt.Sprintf("%s/temperature?cep=%s", os.Getenv("SERVICE_B_ENDPOINT"), cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request")
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot make request: %s", err.Error())
	}
	defer resp.Body.Close()

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %s", err.Error())
	}

	return responseBodyBytes, nil
}
