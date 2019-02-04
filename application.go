package main

import (
	"log"
	"net/http"
	"os"

	"contactBook/conf"
	"contactBook/server"
	"contactBook/handlers"
)

var (
	ServiceAddr = conf.Get().Server
)

func main() {
	logger := log.New(os.Stdout, "contractBook ", log.LstdFlags|log.Lshortfile)

	h := handlers.NewHandlers(logger)

	r := server.New()
	h.SetupRoutes(r)

	if ServiceAddr == "" {
		ServiceAddr = ":5000"
	}
	logger.Println("server listening on " + ServiceAddr)
	err := http.ListenAndServe(ServiceAddr, r)
	if err != nil {
		logger.Fatalf("server failed to start: %v", err)
	}
	
	// TODO how to take care of it in graceful shutdown - process restart
	// Use better logger 
}
