package main

import (
	"fmt"
	"net/http"

	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define Prometheus metrics
var (
	processedMessages = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "processed_messages_total",
			Help: "Total number of messages processed",
		},
	)
	messageProcessingDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "message_processing_duration_seconds",
			Help:    "Duration of message processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(processedMessages)
	prometheus.MustRegister(messageProcessingDuration)
}

func main() {
	// Start NATS connection
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	// Subscribe to NATS subject
	nc.Subscribe("grafana.sub", func(msg *nats.Msg) {
		timer := prometheus.NewTimer(messageProcessingDuration) // Start timer
		defer timer.ObserveDuration()                           // Record duration

		fmt.Printf("Received Message: %s\n", string(msg.Data))
		response := fmt.Sprintf("Processed: %s\n", "eventResponse")
		msg.Respond([]byte(response))

		processedMessages.Inc() // Increment processed messages counter
	})

	// Expose metrics endpoint for Prometheus
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil) // Listen on port 2112 for Prometheus metrics

	select {}
}
