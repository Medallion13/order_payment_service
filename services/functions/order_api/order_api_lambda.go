package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type CreateOrderRequest struct {
	UserId     string `json:"user_id"`
	Item       string `json:"item"`
	Quantity   int    `json:"quantity"`
	TotalPrice int64  `json:"total_price"`
}

type CreateOrderEvent struct {
	OrderID    string `json:"order_id"`
	TotalPrice int64  `json:"total_price"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var order CreateOrderRequest

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

	body, err := json.Marshal(CreateOrderEvent{OrderID: "ramdomId1", TotalPrice: order.TotalPrice})
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
