package kafkaloc

import (
	modalstructs "adsMetrics/modalStructs"
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

func ProcessAdHover(ctx context.Context) {
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			// Context canceled, exit the loop
			fmt.Println("Stopping ProcessAdHover...")
			wg.Wait() // Ensure all goroutines complete
			return

		default:
			msg, err := adhoverreader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message from adhover: %v", err)
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
				case randomNum < 30:
					data.EventType = "adClick"
					data.EventID = uuid.NewString()
					data.Timestamp = time.Now().Unix()
					sendToKafka(usereventswriter, data)

				case randomNum >= 31 && randomNum <= 80:
					data.EventType = "adSkip"
					data.EventID = uuid.NewString()
					data.Timestamp = time.Now().Unix()
					sendToKafka(usereventswriter, data)
				}
			}(msg)
		}
	}
}
