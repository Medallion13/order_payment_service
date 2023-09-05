package main

import (
	"encoding/json"
	"fmt"
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
	// Create util variables
	var orderRequest utils.CreateOrderRequest

	// Inicialize event bridge and Dynamo clients
	dynamo, err := awsUtils.DynamoClient(Table_name)
	if err != nil {
		return awsUtils.CreateBadResponse("Dynamo error", err)
	}
	bridge, err := awsUtils.EventManager(Event_bus_name)
	if err != nil {
		return awsUtils.CreateBadResponse("Event Bridge Error", err)
	}

	// unserialize the event and save the requiered information in the struc object
	err = json.Unmarshal([]byte(request.Body), &orderRequest)
	if err != nil {
		return awsUtils.CreateBadResponse("Api Request Error", customErr.ErrMarsh)
	}

	// Generate unique ID for the order
	orderId := utils.GenKey(15, orderRequest.UserId, orderRequest.Item,
		fmt.Sprintln(orderRequest.TotalPrice), time.Now().String())

	// Create the event for eventBridge and for response
	event := utils.CreateOrderEvent{
		OrderID:    orderId,
		TotalPrice: orderRequest.TotalPrice,
	}

	// Create and put the item in dynamoDB
	err = dynamo.PutItem(utils.OrderTable{
		OrderID:      orderId,
		UserID:       orderRequest.UserId,
		Item:         orderRequest.Item,
		Quantity:     orderRequest.Quantity,
		TotalPrice:   orderRequest.TotalPrice,
		ReadyForShip: false,
		CreateAt:     time.Now().Format(time.RFC822),
	})
	if err != nil {
		return awsUtils.CreateBadResponse("Dynamo error", err)
	}

	// Send the event
	err = bridge.SendEvent("custom.OrderApiFunction", "OrderCreated", event)
	if err != nil {
		return awsUtils.CreateBadResponse("Event Bridge Error", err)
	}

	// Create the response to make the return to api
	response := awsUtils.CreateGoodResponse(event)

	// Send the response to the API
	return response, nil
}

func main() {
	lambda.Start(handler)
}
