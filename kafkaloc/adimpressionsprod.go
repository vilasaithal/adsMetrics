package kafkaloc

import (
	"adsMetrics/generator"
	modalstructs "adsMetrics/models"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

var devices = []string{"mobile", "desktop"}

// needs to push messages to adimpressions topic and userevents topic.

func CreateAdImpressions() int {
	var wg sync.WaitGroup
	errChan := make(chan error)
	counterChan := make(chan int)
	userEventsCounter := 0
	numMessages := 5
	for i := 0; i < numMessages; i++ {
		wg.Add(1)

		go func(messageNum int) {
			defer wg.Done()

			// Generate Campaign and User data
			campaign := generator.CreateCampaign()
			user := generator.CreateUser()

			// Combine them into a single struct
			data := modalstructs.CombinedData{
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
			err = adimpressionswriter.WriteMessages(context.Background(), message)
			if err != nil {
				fmt.Printf("Failed to write message to adimpressions: %v\n", err)
				errChan <- err
			} else {
				fmt.Printf("Message %d sent successfully to adimpressions\n", i)
			}
			err = usereventswriter.WriteMessages(context.Background(), message)
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
