package kafkaloc

import (
	"log"

	"github.com/segmentio/kafka-go"
)

// KafkaWriter is a global Kafka writer instance
var usereventswriter *kafka.Writer
var adimpressionswriter *kafka.Writer
var adhoverwriter *kafka.Writer

// InitKafkaWriter initializes the Kafka writer with the given parameters
func InitKafkaWriter() {
	usereventswriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "userevents",
	})
	adimpressionswriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "adimpressions",
	})
	adhoverwriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "adhover",
	})
}

// CloseKafkaWriter closes the Kafka writer
func CloseKafkaWriter() {
	if usereventswriter != nil {
		err := usereventswriter.Close()
		if err != nil {
			log.Fatalf("failed to close Kafka writer: %v", err)
		}
	}
	if adimpressionswriter != nil {
		err := adimpressionswriter.Close()
		if err != nil {
			log.Fatalf("failed to close Kafka writer: %v", err)
		}
	}
	if adhoverwriter != nil {
		err := adhoverwriter.Close()
		if err != nil {
			log.Fatalf("failed to close Kafka writer: %v", err)
		}
	}
}
