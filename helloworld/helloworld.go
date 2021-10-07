package main

import (
	"context"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"github.com/google/uuid"
)

func receive(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {

	log.Printf("Event received. \n%s\n", event)
	data := &HelloWorld{}

	if err := event.DataAs(data); err != nil {
		log.Printf("Error while extracting cloudevent Data: %s\n", err.Error())
		return nil, cloudevents.NewHTTPResult(400, "failed to convert data: %s", err)
	}

	log.Printf("Hello World Message from received event %q", data.Msg)

	newEvent := cloudevents.NewEvent()
	newEvent.SetID(uuid.New().String())
	newEvent.SetSource("serverless/helloworld")
	newEvent.SetType("servlerless.helloworld.hifromknative")

	if err := newEvent.SetData(cloudevents.ApplicationJSON, HiFromKnative{Msg: "Hi from helloworld app"}); err != nil {
		return nil, cloudevents.NewHTTPResult(500, "failed to set response data: %s", err)
	}
	log.Printf("Responding with event \n%s\n", newEvent)

	return &newEvent, nil
}

func main() {

	log.Print("Hello world sample started")

	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Fatal(c.StartReceiver(context.Background(), receive))

}
