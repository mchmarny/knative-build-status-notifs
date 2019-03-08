package main

import (
	"log"
	"net/http"

	"github.com/knative/pkg/cloudevents"
)


func main() {
	m := cloudevents.NewMux()
	err := m.Handle("google.pubsub.topic.publish", handleMessage)
	if err != nil {
		log.Fatalf("Failed to create handler %s", err)
	}
	http.ListenAndServe(":8080", m)
}