package kafkaloc

import (
	modalstructs "adsMetrics/models"
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

func sendToKafka(writer *kafka.Writer, data modalstructs.CombinedData) {
	// Marshal the updated data back to JSON
	msgValue, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling data: %v", err)
		return
	}

	// Write the message to Kafka
	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: msgValue,
	})
	if err != nil {
		log.Printf("Error writing message to Kafka: %v", err)
	}
}
