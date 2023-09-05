package main

import (
	"encoding/json"
	"log"
	"os"

	awsUtils "github.com/NicoCodes13/order_payment_service/internal/aws"
	custErr "github.com/NicoCodes13/order_payment_service/internal/errors"
	utils "github.com/NicoCodes13/order_payment_service/internal/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var Table_name string

func init() {
	Table_name = os.Getenv("TABLE_NAME")
}

func handler(request events.CloudWatchEvent) error {
	// Util variables
	var paymentEvent utils.PaymentTable

	// Information for log
	log.Printf("Received event of type %q\n", request.DetailType)

	// Unmarsh the request to access the data
	err := json.Unmarshal([]byte(request.Detail), &paymentEvent)
	if err != nil {
		return custErr.ErrUnmarsh
	}

	dynamo, err := awsUtils.DynamoClient(Table_name)
	if err != nil {
		return err
	}
	err = dynamo.UpdateInfo("OrderID", utils.OrderTable{
		OrderID:      paymentEvent.OrderID,
		ReadyForShip: true,
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
