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

func ProcessAdHover() {
	ctx := context.Background()

	var wg sync.WaitGroup

	for {
		msg, err := adhoverreader.ReadMessage(ctx)
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
			case randomNum < 30:
				// Change event type to adError and send to userevents topic
				data.EventType = "adClick"
				data.EventID = uuid.NewString()
				data.Timestamp = time.Now().Unix()
				sendToKafka(usereventswriter, data)

			case randomNum >= 31 && randomNum <= 80:
				// Change event type to adHover and send to both adhover and userevents topics
				data.EventType = "adSkip"
				data.EventID = uuid.NewString()
				data.Timestamp = time.Now().Unix()
				sendToKafka(usereventswriter, data)

			default:
				// Do nothing, just continue to the next message
			}
		}(msg)
	}
	wg.Wait()

}
