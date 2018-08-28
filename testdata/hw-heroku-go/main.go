package main

import (
	"fmt"
	"net/http"
	"os"
)

const Version = "1.0.4"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, Version)
	})

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
