package kafkaloc

import (
	"adsMetrics/generator"
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

type combineddata struct {
	CampaignID      int    `json:"campaign_id"`
	CampaignType    string `json:"campaign_type"`
	CampaignContent string `json:"campaign_content"`
	UserID          int    `json:"user_id"`
	Device          string `json:"device"`
	City            string `json:"city"`
	Age             int    `json:"age"`
	Gender          string `json:"gender"`
	EventID         string `json:"event_id"`
	Timestamp       int64  `json:"timestamp"` // Unix timestamp in seconds
	EventType       string `json:"event_type"`
}

// needs to push messages to adimpressions topic and userevents topic.

func CreateAdImpressions() {
	var wg sync.WaitGroup
	var numMessages = 10
	for i := 0; i < numMessages; i++ {
		wg.Add(1)

		go func(messageNum int) {
			defer wg.Done()

			// Generate Campaign and User data
			campaign := generator.CreateCampaign()
			user := generator.CreateUser()

			// Combine them into a single struct
			data := combineddata{
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
				return
			}

			// Create a Kafka message
			message := kafka.Message{
				Key:   []byte(fmt.Sprintf("key-%d", i)),
				Value: jsonData,
			}

			// Send the JSON message to Kafka
			err = adimpressionswriter.WriteMessages(context.Background(), message)
			if err != nil {
				fmt.Printf("Failed to write message: %v\n", err)
			} else {
				fmt.Printf("Message %d sent successfully\n", i)
			}
		}(i)
	}
	wg.Wait()
}
