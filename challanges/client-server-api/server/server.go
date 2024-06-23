package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	SERVER_PORT                 = ":8080"
	MAX_MS_TO_GET_EXCHANGE_RATE = 200
	MAX_MS_TO_PERSIST_DATA      = 10
)

func main() {
	mux := http.NewServeMux()
	newLogger()

	mux.HandleFunc("GET /cotacao", GetCurrencyExchangeRate)

	slog.Info(fmt.Sprintf("Starting server on port %s", SERVER_PORT))
	http.ListenAndServe(SERVER_PORT, LogMiddleware(mux))
}

func newLogger() {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, nil),
	).
		With("service", "server")

	slog.SetDefault(logger)
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request received", "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

type ExchangeRateResponse struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type ExchangeRate struct {
	Price float64 `json:"price"`
}

func GetCurrencyExchangeRate(w http.ResponseWriter, r *http.Request) {
	client := http.Client{
		Timeout: time.Millisecond * MAX_MS_TO_GET_EXCHANGE_RATE,
	}

	req, err := http.NewRequestWithContext(r.Context(), "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		slog.Error("Error creating request", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		slog.Error("Error getting exchange rate", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	var response ExchangeRateResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		slog.Error("Error decoding response", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	price, err := strconv.ParseFloat(response.USDBRL.Bid, 64)
	if err != nil {
		slog.Error("Error parsing price", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	saveCtx, saveCancel := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond*MAX_MS_TO_PERSIST_DATA))
	defer saveCancel()

	exchangeRate := ExchangeRate{price}
	err = SaveCurrencyExchangeRate(saveCtx, exchangeRate)
	if err != nil {
		slog.Error("Error saving exchange rate", "error", err.Error())
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(exchangeRate)
	if err != nil {
		slog.Error("Error encoding response", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func SaveCurrencyExchangeRate(ctx context.Context, exchangeRate ExchangeRate) error {
	db, err := sql.Open("sqlite3", "./exchange_rate.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS exchange_rate (timestamp DATETIME, price REAL)")
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, "INSERT INTO exchange_rate (timestamp, price) VALUES (?, ?)", time.Now(), exchangeRate.Price)
	if err != nil {
		return err
	}

	return nil
}
