package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

const (
	topic         = "adsCampaign"
	brokerAddress = "localhost:9092"
)

func producer(jsonMessage []byte) {

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})

	defer writer.Close()

	// Create a message
	msg := kafka.Message{
		Value: jsonMessage,
	}

	// Send the message
	err := writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Fatalf("Failed to write message: %v", err)
	}

	fmt.Println("Message sent successfully!")

}
