package aws

import (
	"encoding/json"
	"net/http"

	data "github.com/NicoCodes13/order_payment_service/internal/utils"
	"github.com/aws/aws-lambda-go/events"
)

func CreateGoodResponse(message interface{}) events.APIGatewayProxyResponse {
	// Transform message into json to send in the response of the API
	msg, err := json.Marshal(message)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}
	}

	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Methods": "POST",
		},
	}
	response.StatusCode = http.StatusOK
	response.Body = string(msg)
	return response
}

func CreateBadResponse(err_name string, err error) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Methods": "POST",
		},
	}
	body, err := json.Marshal(data.ErrorApiResponse{ErrorName: err_name, Message: err.Error()})
	if err != nil {
		return response, err
	}

	response.StatusCode = http.StatusInternalServerError
	response.Body = string(body)
	return response, nil
}
