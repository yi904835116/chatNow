package main

import (
	"log"
	"net/http"
	"os"

	"github.com/info344-s18/challenges-yi904835116/servers/gateway/handlers"
)

//main is the main entry point for the server
func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":80"
	}

	mux := http.NewServeMux()

	// mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)

	log.Printf("server is listening at http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
