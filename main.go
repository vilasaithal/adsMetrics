package main

import (
	"adsMetrics/generatorserver"
	"adsMetrics/kafkaloc"
	"os"
	"os/signal"
	"syscall"

	"fmt"
	"log"
	"net/http"
)

func main() {
	kafkaloc.InitKafka()

	go func() {
		kafkaloc.ProcessAdImpressions()
	}()
	go func() {
		kafkaloc.ProcessAdHover()
	}()

	http.HandleFunc("/generate", generatorserver.GenerateHandler)
	server := &http.Server{Addr: ":8080"}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("Server is running on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()
	<-stopChan

	// Gracefully shut down the server
	fmt.Println("Shutting down server...")
	if err := server.Close(); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	// Close Kafka connections
	kafkaloc.CloseKafka()
	fmt.Println("Server and Kafka connections closed")
}
