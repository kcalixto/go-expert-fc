package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	source        string
	executionTime time.Duration
	response      *http.Response
	err           error
}

func main() {
	cep := "01001000"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	responseChan := make(chan Result)

	go getViacep(ctx, cep, responseChan)
	go getBrasilCepAPI(ctx, cep, responseChan)

	mutex := &sync.Mutex{}

	for result := range responseChan {
		mutex.Lock()
		if result.err != nil || result.response == nil {
			fmt.Printf("Erro ao buscar no %s\n", result.source)
			mutex.Unlock()
			continue
		}

		responseBytes, err := io.ReadAll(result.response.Body)
		if err != nil {
			fmt.Printf("Erro ao ler resposta do %s\n", result.source)
			mutex.Unlock()
			return
		}

		fmt.Printf("%s: {\n \tresponseBody: %s\n \texecutionTime: %dms\n \terr: %s\n }", result.source, string(responseBytes), (result.executionTime / time.Millisecond), result.err)
		cancel()
		mutex.Unlock()
		return
	}
}

func getViacep(ctx context.Context, cep string, responseChan chan Result) {
	time.Sleep(randomMs()) // simulates a random delay

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://viacep.com.br/ws/%s/json", cep), nil)
	if err != nil {
		responseChan <- Result{"viacep", time.Since(start), nil, err}
		return
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	responseChan <- Result{"viacep", time.Since(start), resp, err}
}

func getBrasilCepAPI(ctx context.Context, cep string, responseChan chan Result) {
	time.Sleep(randomMs()) // simulates a random delay

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep), nil)
	if err != nil {
		responseChan <- Result{"brasilapi", time.Since(start), nil, err}
		return
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	responseChan <- Result{"brasilapi", time.Since(start), resp, err}
}

func randomMs() time.Duration {
	return time.Duration(rand.Float64()*1000) * time.Millisecond
}
