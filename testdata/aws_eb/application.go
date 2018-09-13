package main

import (
	"fmt"
	"net/http"
	"os"
)

const Version = "1.0.2"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, Version)
	})

	http.ListenAndServe(":"+port, nil)
}
