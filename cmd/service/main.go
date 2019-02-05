package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/mchmarny/knative-build-status-notifs/pkg/build"
	"github.com/mchmarny/knative-build-status-notifs/pkg/pushover"

	"cloud.google.com/go/pubsub"
	"github.com/knative/pkg/cloudevents"
)




func handleMessage(ctx context.Context, msg *pubsub.Message) error {

	ec := cloudevents.FromContext(ctx)
	if ec != nil {
		log.Printf("Received Cloud Event Context as: %v", ec)
	} else {
		log.Printf("No Cloud Event Context found")
	}
	if len(msg.Data) > 0 {
		obj := &build.CloudBuildNotification{}
		err := json.Unmarshal(msg.Data, obj)
		if err != nil {
			log.Printf("Failed to umarshal object notification data: %s\n data was %q", err, string(msg.Data))
			return err
		}
		log.Printf("object notification metadata is: %+v", obj)
		err = pushover.Send(obj)
		if err != nil {
			log.Printf("Failed to send notification %v", err)
			return err
		}
	} else {
		log.Printf("Object Notification event data is empty")
	}

	return nil
}

func main() {
	m := cloudevents.NewMux()
	err := m.Handle("google.pubsub.topic.publish", handleMessage)
	if err != nil {
		log.Fatalf("Failed to create handler %s", err)
	}
	http.ListenAndServe(":8080", m)
}