package main

import (
	"adsMetrics/models"
	"adsMetrics/producers"
	"encoding/json"
	"fmt"
	"net/http"
)

func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count := producers.GetKafkaProducerInstance().CreateAdImpressions()

	fmt.Printf("Generated %d AdImpressions", count)
	resp := models.UserEventsGenerateResp{
		BaseResp: models.BaseResp{
			StatusCode: 200,
			Message:    "Ok",
		},
		UserEventsCount: int32(count),
	}

	message, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(message))

}
