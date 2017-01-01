package main

import (
	"net/http"

	"github.com/christoph-k/go-http-logger"
)

func serve(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!\n"))
}

func main() {
	http.Handle("/", logger.New(http.HandlerFunc(serve)))
	http.ListenAndServe(":8080", nil)
}
