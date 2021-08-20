package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/go-github/v38/github"
	"github.com/sethvargo/go-retry"
)

type LambdaHandler func(context.Context) error

type Event struct {
	Type   string
	Detail json.RawMessage
}

func EventsHandler(githubClient *github.Client, httpClient *http.Client, apiEndpoint string) LambdaHandler {
	return func(ctx context.Context) error {
		log.Printf("Using API Endpoint: %s", apiEndpoint)

		ctx, cancelFn := context.WithTimeout(ctx, time.Second*5)
		defer cancelFn()

		// more advanced usage would check the rate limit from the response and add jitter
		events, _, err := githubClient.Activity.ListEvents(ctx, nil)
		if err != nil {
			log.Printf("Error listing events: %v", err)
			return err
		}

		log.Printf("Got %d events", len(events))

		for _, e := range events {
			body := Event{
				Type:   e.GetType(),
				Detail: e.GetRawPayload(),
			}

			payload, err := json.Marshal(body)
			if err != nil {
				log.Printf("Unable to marshal message: %v", e.GetRawPayload())
				continue
			}

			// a further enhancement would be to batch these to be posted
			err = PostEvent(ctx, httpClient, apiEndpoint, string(payload))
			if err != nil {
				// we could push to a DLQ here
				log.Printf("Unable to post payload (%s) to %s", e.GetRawPayload(), apiEndpoint)
			}
		}

		return nil
	}
}

func PostEvent(ctx context.Context, httpClient *http.Client, apiEndpoint string, event string) error {
	backoff, err := retry.NewFibonacci(time.Second * 1)
	if err != nil {
		return err
	}

	err = retry.Do(ctx, retry.WithJitterPercent(5, backoff), func(_ context.Context) error {
		buf := bytes.NewBuffer([]byte(event))

		_, err := httpClient.Post(apiEndpoint, "application/json", buf)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func main() {
	// initialise a new github client to pass to the handler
	client := github.NewClient(nil)

	// get env var for the events api we'll post events to
	apiEndpoint := os.Getenv("EVENTS_ENDPOINT")
	if len(apiEndpoint) == 0 {
		log.Fatalln("Endpoint environment variable (EVENTS_ENDPOINT) not set")
	}

	lambda.Start(EventsHandler(client, http.DefaultClient, apiEndpoint))
}
