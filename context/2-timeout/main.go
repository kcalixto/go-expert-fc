package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, time.Second*3)
		fmt.Println("started")
		defer cancel()

		res := make(chan []byte)

		go func() {
			time.Sleep(time.Second * 10)
			res <- []byte("hello!!")
		}()

		select {
		case <-ctx.Done():
			message := fmt.Sprintf("request canceled: %s", ctx.Err().Error())
			fmt.Println(message)

			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte(message))
		case result := <-res:
			w.WriteHeader(http.StatusOK)
			w.Write(result)
		}
	})

	http.ListenAndServe(":8080", nil)
}
