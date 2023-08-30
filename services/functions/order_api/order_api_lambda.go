package main

import (
	"encoding/json"

	awsUtils "github.com/NicoCodes13/order_payment_service/internal/aws"
	data "github.com/NicoCodes13/order_payment_service/internal/utils"
	internalErrors "github.com/NicoCodes13/order_payment_service/internal/errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var order data.CreateOrderEvent
	// unserialize the event and save the requiered information in the struc object
	err := json.Unmarshal([]byte(request.Body), &order)
	err = nil
	if err != nil {
		return awsUtils.CreateBadResponse("Problem unmarshal the request from the api")
	}

	body, err := json.Marshal(data.CreateOrderEvent{OrderID: "ramdomId1", TotalPrice: order.TotalPrice})
	if err != nil {
		return awsUtils.CreateBadResponse("Problem marshal body response")
	}

	response = awsUtils.CreateGoodResponse(string(body))

	return response, nil
}

func main() {
	lambda.Start(handler)
}
