package main

import (
	generatorserver "adsMetrics/GeneratorServer"
	"adsMetrics/kafkaloc"

	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	kafkaloc.InitKafkaWriter()
	http.HandleFunc("/generate", generatorserver.GenerateHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
