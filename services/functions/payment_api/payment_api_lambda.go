package main

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	awsUtils "github.com/NicoCodes13/order_payment_service/internal/aws"
	customErr "github.com/NicoCodes13/order_payment_service/internal/errors"
	"github.com/NicoCodes13/order_payment_service/internal/utils"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var payment_request utils.ProcessPaymentData

	// unserialize the event and save the requiered information in the struc object
	err := json.Unmarshal([]byte(request.Body), &payment_request)
	if err != nil {
		return awsUtils.CreateBadResponse("Api Request Error", customErr.ErrMarsh)
	}
	payment_request.Status = validation(payment_request.Status)

	// Create the body of the response using the struct CreateOrderEvent
	body, err := json.Marshal(payment_request)
	if err != nil {
		return awsUtils.CreateBadResponse("API body response Error", customErr.ErrMarsh)
	}

	response := awsUtils.CreateGoodResponse(string(body))

	return response, nil
}

func main() {
	lambda.Start(handler)
}

func validation(status string) string {
	if strings.ToLower(status) == "paid" {
		return "complete"
	}
	return "incomplete"

}
