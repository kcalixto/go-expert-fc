package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server := &http.Server{Addr: ":8080"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 3)
		w.Write([]byte("hello!!"))
	})

	go func() {
		log.Println("server started")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	// Operational system signal
	stop := make(chan os.Signal, 1)
	// Notify will call the stop channel when a signal of any of the following types is received
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// awaits for a signal to be received
	select {
	case <-stop:
		// When this context achieve it's timeout the server will be shutdown forcefully
		// If the server shutdown before 30s it just works gracefully
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		log.Println("shutting down server...")
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("could not gracefully shutdown the server: %v\n", err)
		}

		log.Println("server stoped")
	}
}
