package send

import (
	"bytes"
	"encoding/json"
	"github.com/gorlug/metrics-backend/metrics"
	"log"
	"net/http"
	"os"
	"time"
)

func SendMetric(metricBuilder *metrics.MetricBuilder, url string) {
	metricBuilder.
		WithHost(getHostname()).WithTimestamp(time.Now())
	metric := metricBuilder.Build()
	log.Printf("sending Metric %v\n", metric.String())

	metricJson, marshalErr := json.Marshal(metric)

	if marshalErr != nil {
		log.Fatal("Could not marshal metric into JSON")
	}

	// Make request with marshalled JSON as the POST body
	response, err := http.Post(url, "application/json",
		bytes.NewBuffer(metricJson))

	if err != nil {
		log.Fatal("Could not make POST request")
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("Error Response status: %v\n", response.Status)
	}
}

func getHostname() string {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return name
}
