package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	MAX_MS_TO_GET_QUOTATION = 300
	QUOTATION_FILE_NAME     = "cotacao.txt"
)

func main() {
	newLogger()

	quotation, err := GetQuotation()
	if err != nil {
		slog.Error("Error getting quotation", "error", err)
		return
	}

	err = Save(quotation)
	if err != nil {
		slog.Error("Error saving quotation", "error", err)
		return
	}
}

func newLogger() {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, nil),
	).
		With("service", "client")

	slog.SetDefault(logger)
}

type QuotationResponse struct {
	Price float64 `json:"price"`
}

// Gets the quotation from the server
func GetQuotation() (quotation QuotationResponse, err error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond*MAX_MS_TO_GET_QUOTATION))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Error getting exchange rate", "error", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return quotation, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&quotation); err != nil {
		slog.Error("Error decoding response", "error", err)
		return quotation, err
	}

	return quotation, nil
}

// Saves the quotation to a file
func Save(quotation QuotationResponse) error {
	currentTime := time.Now().Format(time.RFC3339)
	newLine := fmt.Sprintf("%s Dólar: %.2f\n", currentTime, quotation.Price)

	oldLines, err := os.ReadFile(QUOTATION_FILE_NAME)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	fileContent := fmt.Sprintf("%s%s", string(oldLines), newLine)

	err = os.WriteFile(QUOTATION_FILE_NAME, []byte(fileContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
