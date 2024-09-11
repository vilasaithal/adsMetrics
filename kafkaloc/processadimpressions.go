package kafkaloc

import (
	modalstructs "adsMetrics/models"
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

func ProcessAdImpressions(ctx context.Context) {
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			// Context canceled, exit the loop
			fmt.Println("Stopping ProcessAdImpressions...")
			wg.Wait() // Ensure all goroutines complete
			return

		default:
			msg, err := adimpressionsreader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message from adimpressions: %v", err)
				continue
			}

			wg.Add(1)

			go func(msg kafka.Message) {
				defer wg.Done()

				var data modalstructs.CombinedData
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
					sendToKafka(usereventswriter, data)

				case randomNum >= 6 && randomNum <= 40:
					data.EventType = "adHover"
					data.EventID = uuid.NewString()
					data.Timestamp = time.Now().Unix()
					sendToKafka(adhoverwriter, data)
					sendToKafka(usereventswriter, data)
				}
			}(msg)
		}
	}
}
