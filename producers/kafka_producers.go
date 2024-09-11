package producers

import (
	"adsMetrics/generator"
	"adsMetrics/models"
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

type KafkaProducer struct {
	Usereventswriter    *kafka.Writer
	Adimpressionswriter *kafka.Writer
	Adhoverwriter       *kafka.Writer
}

var (
	kafkaProducersInstance *KafkaProducer
	once                   sync.Once // ensures the instance is created only once
)

func GetKafkaProducerInstance() *KafkaProducer {
	once.Do(func() {
		log.Println("Creating producer instance")
		kafkaProducersInstance = &KafkaProducer{}
	})
	return kafkaProducersInstance
}

func (p *KafkaProducer) InitProducers() {
	p.Usereventswriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "userevents",
	})

	p.Adimpressionswriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "adimpressions",
	})

	p.Adhoverwriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "adhover",
	})

}

func (p *KafkaProducer) CloserProducers() error {
	if p.Usereventswriter != nil {
		err := p.Usereventswriter.Close()
		if err != nil {
			log.Fatalf("failed to close Kafka writer: %v", err)
			return err
		}
	}
	if p.Adimpressionswriter != nil {
		err := p.Adimpressionswriter.Close()
		if err != nil {
			log.Fatalf("failed to close Kafka writer: %v", err)
			return err
		}
	}
	if p.Adhoverwriter != nil {
		err := p.Adhoverwriter.Close()
		if err != nil {
			log.Fatalf("failed to close Kafka writer: %v", err)
			return err
		}
	}

	return nil
}

func (p *KafkaProducer) GenerateUserEvents(data models.CombinedData) error {
	return p.KafkaEventWrite(p.Usereventswriter, data)
}

func (p *KafkaProducer) Generateadimpressions(data models.CombinedData) error {
	return p.KafkaEventWrite(p.Adimpressionswriter, data)
}

func (p *KafkaProducer) Generateadhover(data models.CombinedData) error {
	return p.KafkaEventWrite(p.Adhoverwriter, data)
}

func (p *KafkaProducer) KafkaEventWrite(writer *kafka.Writer, data models.CombinedData) error {
	// Marshal the updated data back to JSON
	msgValue, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling data: %v", err)
		return err
	}

	// Write the message to Kafka
	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: msgValue,
	})
	if err != nil {
		log.Printf("Error writing message to Kafka: %v", err)
		return err
	}

	return nil
}

func (p *KafkaProducer) StartAdImpressionProducer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	count := 0

	for {
		select {
		case <-ticker.C:
			count += p.CreateAdImpressions()
		case <-ctx.Done(): // Context cancellation (client disconnect)
			fmt.Printf("StartAdImpressionProducer received shutdown signal. Exiting... \n")
			return
		}
	}
}

func (p *KafkaProducer) CreateAdImpressions() int {
	var wg sync.WaitGroup
	errChan := make(chan error)
	counterChan := make(chan int)
	userEventsCounter := 0
	numMessages := 5
	devices := []string{"mobile", "desktop"}
	for i := 0; i < numMessages; i++ {
		wg.Add(1)

		go func(messageNum int) {
			defer wg.Done()

			// Generate Campaign and User data
			campaign := generator.CreateCampaign()
			user := generator.CreateUser()

			// Combine them into a single struct
			data := models.CombinedData{
				CampaignID:      campaign.CampaignID,
				CampaignType:    campaign.CampaignType,
				CampaignContent: campaign.CampaignContent,
				UserID:          user.UserID,
				Device:          devices[rand.Intn(len(devices))],
				City:            user.City,
				Age:             user.Age,
				Gender:          user.Gender,
				EventID:         uuid.NewString(),
				Timestamp:       time.Now().Unix(), // Get the current Unix timestamp
				EventType:       "adImpression",
			}

			// Marshal the combined struct into JSON
			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Printf("Failed to marshal data: %v\n", err)
				errChan <- err
			}

			// Create a Kafka message
			message := kafka.Message{
				Key:   []byte(fmt.Sprintf("key-%d", i)),
				Value: jsonData,
			}

			// Send the JSON message to Kafka
			err = p.Adimpressionswriter.WriteMessages(context.Background(), message)
			if err != nil {
				fmt.Printf("Failed to write message to adimpressions: %v\n", err)
				errChan <- err
			} else {
				fmt.Printf("Message %d sent successfully to adimpressions\n", i)
			}
			err = p.Usereventswriter.WriteMessages(context.Background(), message)
			if err != nil {
				fmt.Printf("Failed to write message to user events: %v\n", err)
				errChan <- err
			} else {
				fmt.Printf("Message %d sent successfully to user events\n", i)
				counterChan <- 1
			}
		}(i)
	}
	go func() {
		wg.Wait()
		close(counterChan)
	}()

	for count := range counterChan {
		userEventsCounter += count
	}

	fmt.Println("Finished creating and sending messages")

	return userEventsCounter
}
