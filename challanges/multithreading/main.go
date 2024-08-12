package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	source   string
	response *http.Response
	err      error
}

func main() {
	cep := "01001000"

	ctx, cancel := context.WithCancel(context.Background())
	responseChan := make(chan Result)

	go getViacep(ctx, cep, responseChan)
	go getBrasilCepAPI(ctx, cep, responseChan)

	mutex := &sync.Mutex{}

	for {
		select {
		case result := <-responseChan:
			mutex.Lock()
			if result.err != nil || result.response == nil {
				fmt.Printf("Erro ao buscar no %s\n", result.source)
				mutex.Unlock()
				continue
			}

			var body map[string]interface{}
			err := json.NewDecoder(result.response.Body).Decode(&body)
			if err != nil {
				fmt.Printf("Erro ao decodificar resposta do %s\n", result.source)
				mutex.Unlock()
				return
			}
			fmt.Printf("%s: %s\n", result.source, body)
			cancel()
			mutex.Unlock()
			return
		}
	}
}

func getViacep(ctx context.Context, cep string, responseChan chan Result) {
	time.Sleep(randomMs())
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://viacep.com.br/ws/%s/json", cep), nil)
	if err != nil {
		responseChan <- Result{"viacep", nil, err}
		return
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	responseChan <- Result{"viacep", resp, err}
}

func getBrasilCepAPI(ctx context.Context, cep string, responseChan chan Result) {
	time.Sleep(randomMs())
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep), nil)
	if err != nil {
		responseChan <- Result{"brasilapi", nil, err}
		return
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	responseChan <- Result{"brasilapi", resp, err}
}

func randomMs() time.Duration {
	return time.Duration(rand.Float64()*1000) * time.Millisecond
}
