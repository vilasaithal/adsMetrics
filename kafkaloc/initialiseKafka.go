package kafkaloc

import (
	"log"

	"github.com/segmentio/kafka-go"
)

// KafkaWriter is a global Kafka writer instance
var usereventswriter *kafka.Writer
var adimpressionswriter *kafka.Writer
var adhoverwriter *kafka.Writer
var adimpressionsreader *kafka.Reader
var adhoverreader *kafka.Reader

// InitKafkaWriter initializes the Kafka writer with the given parameters
func InitKafka() {
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
	adimpressionsreader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "adimpressions",
		GroupID: "groupA", // Set a consumer group ID
	})
	adhoverreader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "adhover",
		GroupID: "groupA", // Set a consumer group ID
	})
}

// CloseKafkaWriter closes the Kafka writer
func CloseKafka() {
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
	if adimpressionsreader != nil {
		err := adimpressionsreader.Close()
		if err != nil {
			log.Fatalf("failed to close Kafka adimpressionreader: %v", err)
		}
	}
	if adhoverreader != nil {
		err := adhoverreader.Close()
		if err != nil {
			log.Fatalf("failed to close Kafka adhoverreader: %v", err)
		}
	}
}
