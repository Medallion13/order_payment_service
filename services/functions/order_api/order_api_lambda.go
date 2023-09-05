package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	awsUtils "github.com/NicoCodes13/order_payment_service/internal/aws"
	customErr "github.com/NicoCodes13/order_payment_service/internal/errors"
	utils "github.com/NicoCodes13/order_payment_service/internal/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var Table_name string
var Event_bus_name string

func init() {
	Table_name = os.Getenv("TABLE_ORDER")
	Event_bus_name = os.Getenv("EVENT_BUS_NAME")
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var order_request utils.CreateOrderRequest

	// unserialize the event and save the requiered information in the struc object
	err := json.Unmarshal([]byte(request.Body), &order_request)
	if err != nil {
		return awsUtils.CreateBadResponse("Api Request Error", customErr.ErrMarsh)
	}

	// Generate unique ID for the order
	orderId := utils.GenKey(15, order_request.UserId, order_request.Item,
		fmt.Sprintln(order_request.TotalPrice), time.Now().String())

	// Create the event for eventBridge
	event := utils.CreateOrderEvent{
		OrderID:    orderId,
		TotalPrice: order_request.TotalPrice,
	}

	log.Println(Event_bus_name)
	// Send the event
	bridge, err := awsUtils.EventManager(Event_bus_name)
	if err != nil {
		awsUtils.CreateBadResponse("Event Bridge Error", err)
	}

	err = bridge.SendEvent("custom.OrderApiFunction", "create_payment", event)
	if err != nil {
		awsUtils.CreateBadResponse("Event Bridge Error", err)
	}

	// Create the response to make the return to api
	response := awsUtils.CreateGoodResponse(event)

	// Send the response to the API
	return response, nil
}

func main() {
	lambda.Start(handler)
}
