package generatorserver

import (
	"adsMetrics/kafkaloc"
	"fmt"
	"net/http"
	"time"
)

func Generateadimpressions(w http.ResponseWriter, r *http.Request) {

	// Use the request's context to handle cancellation
	ctx := r.Context()
	stop := make(chan struct{})

	// Start a goroutine to continuously generate messages
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				kafkaloc.CreateAdImpressions()
			case <-stop:
				return
			case <-ctx.Done(): // Context cancellation (client disconnect)
				close(stop)
				return
			}
		}
	}()

	// Wait for the client to disconnect or the context to be canceled
	<-ctx.Done()

	fmt.Fprintf(w, "Ad Message generation stopped")
}
