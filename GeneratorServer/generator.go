package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
)

type Campaign struct {
	OwnerID  int    `json:"owner_id"`
	AdID     int    `json:"ad_id"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	Interest string `json:"interest"`
	City     string `json:"city"`
	State    string `json:"state"`
}

var (
	interests = []string{"Sports", "Technology", "Fashion", "Food", "Travel", "Music", "Movies", "Books", "Art", "Fitness"}
	cities    = []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose"}
	states    = []string{"NY", "CA", "IL", "TX", "AZ", "PA", "FL", "OH", "MI", "GA"}
	genders   = []string{"Male", "Female", "Other"}
)

func generateCampaign() Campaign {
	return Campaign{
		OwnerID:  rand.Intn(10000) + 1,
		AdID:     rand.Intn(50) + 1,
		Age:      rand.Intn(43) + 18, // 18 to 60
		Gender:   genders[rand.Intn(len(genders))],
		Interest: interests[rand.Intn(len(interests))],
		City:     cities[rand.Intn(len(cities))],
		State:    states[rand.Intn(len(states))],
	}
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	campaign := generateCampaign()

	// Convert the campaign to JSON
	campaignJSON, err := json.Marshal(campaign)
	if err != nil {
		log.Fatalf("Failed to marshal campaign to JSON: %v", err)
	}

	// Send the JSON message to the Kafka topic
	producer(campaignJSON)
}
