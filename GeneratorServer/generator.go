package generatorserver

import (
	"adsMetrics/kafkaloc"
	"fmt"
	"net/http"
	"time"
)

func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	stop := make(chan struct{})

	// Start a goroutine to continuously generate messages
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				kafkaloc.CreateAdImpressions()
			case <-stop:
				return
			}
		}
	}()

	// Wait for the client to disconnect
	notify := w.(http.CloseNotifier).CloseNotify()
	<-notify

	// Signal the goroutine to stop
	close(stop)

	fmt.Fprintf(w, "Message generation stopped")
}
