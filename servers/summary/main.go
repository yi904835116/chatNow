package main

import (
	"log"
	"net/http"
	"os"

	"github.com/info344-s18/challenges-yi904835116/servers/summary/handlers"
)

func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":80"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)

	log.Printf("Server is listening at http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
