package consumers

import (
	"adsMetrics/models"
	"adsMetrics/producers"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Adimpressionsreader *kafka.Reader
	Adhoverreader       *kafka.Reader
}

var (
	kafkaConsumersInstance *KafkaConsumer
	once                   sync.Once // ensures the instance is created only once
)

func GetKafkaConsumerInstance() *KafkaConsumer {
	once.Do(func() {
		log.Println("Creating producer instance")
		kafkaConsumersInstance = &KafkaConsumer{}
	})
	return kafkaConsumersInstance
}

func (c *KafkaConsumer) InitConsumers() {
	c.Adimpressionsreader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "adimpressions",
		GroupID: "adimpression_cons_group", // Set a consumer group ID
	})

	c.Adhoverreader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "adhover",
		GroupID: "adhover_cons_group", // Set a consumer group ID
	})

}

func (c *KafkaConsumer) CloseConsumers() error {
	if c.Adimpressionsreader != nil {
		err := c.Adimpressionsreader.Close()
		if err != nil {
			log.Fatalf("failed to close Kafka adimpressionreader: %v", err)
			return err
		}
	}

	if c.Adhoverreader != nil {
		err := c.Adhoverreader.Close()
		if err != nil {
			log.Fatalf("failed to close Kafka adhoverreader: %v", err)
			return err
		}
	}

	return nil
}

func (c *KafkaConsumer) StartAdImpressionConsumer(ctx context.Context, wg *sync.WaitGroup, p *producers.KafkaProducer) {

	for {
		select {
		case <-ctx.Done():
			// Context canceled, exit the loop
			fmt.Println("Stopping ProcessAdImpressions...")
			return

		default:
			msg, err := c.Adimpressionsreader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message from adimpressions: %v", err)
				continue
			}

			var data models.CombinedData
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				return
			}

			randomNum := rand.Intn(100) + 1

			switch {
			case randomNum < 5:
				data.EventType = "adError"
				data.EventID = uuid.NewString()
				data.Timestamp = time.Now().Unix()
				p.GenerateUserEvents(data)

			case randomNum >= 6 && randomNum <= 40:
				data.EventType = "adHover"
				data.EventID = uuid.NewString()
				data.Timestamp = time.Now().Unix()
				p.Generateadhover(data)
				p.GenerateUserEvents(data)
			}
		}
	}
}

func (c *KafkaConsumer) StartAdHoverConsumer(ctx context.Context, wg *sync.WaitGroup, p *producers.KafkaProducer) {
	for {
		select {
		case <-ctx.Done():
			// Context canceled, exit the loop
			fmt.Println("Stopping ProcessAdHover...")
			return

		default:
			msg, err := c.Adhoverreader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message from adhover: %v", err)
				continue
			}

			var data models.CombinedData
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				return
			}

			randomNum := rand.Intn(100) + 1

			switch {
			case randomNum < 30:
				data.EventType = "adClick"
				data.EventID = uuid.NewString()
				data.Timestamp = time.Now().Unix()
				p.GenerateUserEvents(data)

			case randomNum >= 31 && randomNum <= 80:
				data.EventType = "adSkip"
				data.EventID = uuid.NewString()
				data.Timestamp = time.Now().Unix()
				p.GenerateUserEvents(data)
			}
		}
	}
}
