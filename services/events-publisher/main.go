package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

type LambdaHandler func(context.Context, *events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error)

type Event struct {
	Type   string
	Detail json.RawMessage
}

func EventsPublisher(eb *eventbridge.Client, busName string) LambdaHandler {
	return func(ctx context.Context, req *events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
		log.Printf("received: %s", req.Body)

		payload := Event{}
		if err := json.Unmarshal([]byte(req.Body), &payload); err != nil {
			return &events.APIGatewayV2HTTPResponse{
				StatusCode: 400,
			}, nil
		}

		eb.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					EventBusName: aws.String(busName),
					Source:       aws.String("com.amazonaws.devax.demo"),
					DetailType:   aws.String(payload.Type),
					Detail:       aws.String(string(payload.Detail)),
				},
			},
		})

		return &events.APIGatewayV2HTTPResponse{
			StatusCode: 201,
			Body:       req.Body,
		}, nil
	}
}

func main() {
	busName := os.Getenv("EVENTBUS_NAME")
	if len(busName) == 0 {
		log.Fatalln("Event bus name environment variable (EVENTBUS_NAME) not set")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Unable to load config: %v", err)
	}

	svc := eventbridge.NewFromConfig(cfg)

	lambda.Start(EventsPublisher(svc, busName))
}
