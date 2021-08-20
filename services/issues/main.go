package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

type LambdaHandler func(context.Context) error

func IssuesHandler() LambdaHandler {
	return func(ctx context.Context) error {
		log.Printf("Matched against a rule")
		return nil
	}
}

func main() {
	lambda.Start(IssuesHandler())
}
