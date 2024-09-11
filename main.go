package main

import (
	"adsMetrics/consumers"
	"adsMetrics/producers"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"fmt"
	"log"
	"net/http"
)

func main() {
	kafkaProducers := producers.GetKafkaProducerInstance()
	kafkaConsumers := consumers.GetKafkaConsumerInstance()

	kafkaProducers.InitProducers()
	kafkaConsumers.InitConsumers()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	wg.Add(1)
	go kafkaConsumers.StartAdImpressionConsumer(ctx, &wg, kafkaProducers)

	wg.Add(1)
	go kafkaConsumers.StartAdHoverConsumer(ctx, &wg, kafkaProducers)

	// TODO : uncomment when ready
	wg.Add(1)
	go kafkaProducers.StartAdImpressionProducer(ctx, &wg)

	http.HandleFunc("/generate", GenerateHandler)
	server := &http.Server{Addr: ":8080"}

	// go producers.StartEvents(stopChan)

	go func() {
		fmt.Println("Server is running on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-signalChan
	fmt.Printf("\nReceived signal: %s. Shutting down...\n", sig)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	cancel() // Signal the goroutines to stop

	// Wait for Kafka consumer and producer to finish
	wg.Wait()

	// Close Kafka consumer and producer
	if err := kafkaConsumers.CloseConsumers(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)
	}
	if err := kafkaProducers.CloserProducers(); err != nil {
		log.Printf("Error closing Kafka producer: %v", err)
	}

	fmt.Println("Shutdown complete.")

}
