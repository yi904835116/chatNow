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
		addr = ":443"
	}

	//get the TLS key and cert paths from environment variables
	//this allows us to use a self-signed cert/key during development
	//and the Let's Encrypt cert/key in production
	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")

	if len(tlsKeyPath) == 0 || len(tlsCertPath) == 0 {
		log.Fatal("Please set TLSCERT and TLSKEY")
	}

	// Set up Redis connection.
	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		redisAddr = "localhost:6379"
	}

	mux := http.NewServeMux()

	// mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)

	log.Printf("Server is listening at https://%s\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, mux))
}
