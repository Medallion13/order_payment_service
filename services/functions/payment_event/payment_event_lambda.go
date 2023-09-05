package main

import (
	"encoding/json"
	"fmt"

	"github.com/NicoCodes13/order_payment_service/internal/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davecgh/go-spew/spew"
)

func handler(request events.CloudWatchEvent) error {
	var orderEvent utils.CreateOrderEvent
	fmt.Printf("Received event of type %q\n", request.DetailType)
	err := json.Unmarshal([]byte(request.Detail), &orderEvent)
	if err != nil {
		return err
	}
	spew.Dump(orderEvent)
	return nil
}

func main() {
	lambda.Start(handler)
}
