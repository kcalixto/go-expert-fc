package main

import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
