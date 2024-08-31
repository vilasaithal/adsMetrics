package kafkaloc

import (
	modalstructs "adsMetrics/modalStructs"
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func ProcessAdImpressions() {
	ctx := context.Background()

	var wg sync.WaitGroup

	for {
		msg, err := adimpressionsreader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		// Increment the WaitGroup counter
		wg.Add(1)

		// Process each message in a new goroutine
		go func(msg kafka.Message) {
			defer wg.Done() // Decrement the counter when the goroutine completes

			// Unmarshal the message into CombinedData struct
			var data modalstructs.CombinedData
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				return
			}

			// Generate a random number between 1 and 100
			randomNum := rand.Intn(100) + 1

			// Process based on the random number
			switch {
			case randomNum < 5:
				// Change event type to adError and send to userevents topic
				data.EventType = "adError"
				data.EventID = uuid.NewString()
				data.Timestamp = time.Now().Unix()
				sendToKafka(usereventswriter, data)

			case randomNum >= 6 && randomNum <= 40:
				// Change event type to adHover and send to both adhover and userevents topics
				data.EventType = "adHover"
				data.EventID = uuid.NewString()
				data.Timestamp = time.Now().Unix()
				sendToKafka(adhoverwriter, data)
				sendToKafka(usereventswriter, data)

			default:
				// Do nothing, just continue to the next message
			}
		}(msg)
	}

	// Wait for all goroutines to finish before exiting
	wg.Wait()
}
