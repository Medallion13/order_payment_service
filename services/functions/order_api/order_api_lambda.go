package main

import (
	"encoding/json"
	"fmt"
	"time"

	awsUtils "github.com/NicoCodes13/order_payment_service/internal/aws"
	customErr "github.com/NicoCodes13/order_payment_service/internal/errors"
	utils "github.com/NicoCodes13/order_payment_service/internal/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var order_request utils.CreateOrderRequest

	// unserialize the event and save the requiered information in the struc object
	err := json.Unmarshal([]byte(request.Body), &order_request)
	if err != nil {
		return awsUtils.CreateBadResponse("Api Request Error", customErr.ErrMarsh)
	}

	// Generate unique ID for the order
	orderId := utils.GenKey(15, order_request.UserId, order_request.Item, fmt.Sprintln(order_request.TotalPrice), time.Now().String())

	// Create the body of the response using the struct CreateOrderEvent
	body, err := json.Marshal(utils.CreateOrderEvent{OrderID: orderId, TotalPrice: order_request.TotalPrice})
	if err != nil {
		return awsUtils.CreateBadResponse("API body response Error", customErr.ErrMarsh)
	}

	// Create the response to make the return to api
	response := awsUtils.CreateGoodResponse(string(body))

	return response, nil
}

func main() {
	lambda.Start(handler)
}
