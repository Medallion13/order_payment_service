package aws

import (
	"encoding/json"
	"net/http"

	data "github.com/NicoCodes13/order_payment_service/internal/utils"
	"github.com/aws/aws-lambda-go/events"
)

func CreateGoodResponse(message string) events.APIGatewayProxyResponse {
	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Methods": "POST",
		},
	}
	response.StatusCode = http.StatusOK
	response.Body = message
	return response
}

func CreateBadResponse(err_msg string) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Methods": "POST",
		},
	}
	body, err := json.Marshal(data.ErrorApiResponse{ErrorName: "LambdaError", Message: err_msg})
	if err != nil {
		return response, err
	}

	response.StatusCode = http.StatusBadRequest
	response.Body = string(body)
	return response, nil
}
