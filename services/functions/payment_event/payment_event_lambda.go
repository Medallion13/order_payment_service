package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davecgh/go-spew/spew"
)

func handler(request events.CloudWatchEvent) error {
	fmt.Printf("Received event of type %q\n", request.DetailType)
	spew.Dump(request)
	return nil
}

func main() {
	lambda.Start(handler)
}
