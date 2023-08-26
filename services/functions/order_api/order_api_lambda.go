package main

import (
	"encoding/json"
	"net/http"

	data "github.com/NicoCodes13/order_payment_service/internal/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var order data.CreateOrderEvent

	// Creating a response and adding headers
	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Methods": "POST",
		},
	}

	// unserialize the event and save the requiered information in the struc object
	err := json.Unmarshal([]byte(request.Body), &order)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		return response, err
	}

	body, err := json.Marshal(data.CreateOrderEvent{OrderID: "ramdomId1", TotalPrice: order.TotalPrice})
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		return response, err
	}

	response.StatusCode = http.StatusOK
	response.Body = string(body)
	return response, nil
}

func main() {
	lambda.Start(handler)
}
