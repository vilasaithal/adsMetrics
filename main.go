package main

import (
	"adsMetrics/generatorserver"
	"adsMetrics/kafkaloc"
	"context"
	"os"
	"os/signal"
	"syscall"

	"fmt"
	"log"
	"net/http"
)

func main() {
	kafkaloc.InitKafka()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		kafkaloc.ProcessAdImpressions(ctx)
	}()
	go func() {
		kafkaloc.ProcessAdHover(ctx)
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

	fmt.Println("Shutting down server...")
	if err := server.Close(); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	cancel() // Signal the goroutines to stop
	kafkaloc.CloseKafka()
	fmt.Println("Server and Kafka connections closed")
}
