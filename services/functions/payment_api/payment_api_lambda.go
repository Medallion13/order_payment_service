package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	awsUtils "github.com/NicoCodes13/order_payment_service/internal/aws"
	customErr "github.com/NicoCodes13/order_payment_service/internal/errors"
	"github.com/NicoCodes13/order_payment_service/internal/utils"
)

var Table_name string
var Event_bus_name string

func init() {
	Table_name = os.Getenv("TABLE_ORDER")
	Event_bus_name = os.Getenv("EVENT_BUS_NAME")
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Create util variables
	var payment_request utils.ProcessPaymentData
	var item utils.PaymentTable

	// unserialize the event and save the requiered information in the struc object
	err := json.Unmarshal([]byte(request.Body), &payment_request)
	if err != nil {
		return awsUtils.CreateBadResponse("Api Request Error", customErr.ErrMarsh)
	}
	payment_request.Status = validation(payment_request.Status)

	// Inicialize event bridge and Dynamo clients
	dynamo, err := awsUtils.DynamoClient(Table_name)
	if err != nil {
		return awsUtils.CreateBadResponse("Dynamo error", err)
	}
	bridge, err := awsUtils.EventManager(Event_bus_name)
	if err != nil {
		return awsUtils.CreateBadResponse("Event Bridge Error", err)
	}

	// search for order id created in payment table
	err = dynamo.GetItem("OrderID", payment_request.OrderID, &item)
	if err != nil {
		return awsUtils.CreateBadResponse("Dynamo Error", err)
	}
	// If the item doesn't exist change the status of the payment respons to Not found
	// Else update the status of the payment in the table
	if utils.IsEmpty(item.OrderID) {
		payment_request.Status = "Not Found"
	} else {
		// Update item information for the dynamoDB table
		item.PaymentStatus = payment_request.Status
		// Update payment table
		dynamo.UpdateInfo("OrderID", item)

		// Create event for Order event lambda
		err = bridge.SendEvent("custom.PaymentApiFunction", "PaymentComplete", item)
		if err != nil {
			return awsUtils.CreateBadResponse("Event Bridge Error", err)
		}
	}

	// Create the response
	response := awsUtils.CreateGoodResponse(payment_request)

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
