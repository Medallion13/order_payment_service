package aws

import (
	"context"
	"encoding/json"
	"log"

	customErr "github.com/NicoCodes13/order_payment_service/internal/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/davecgh/go-spew/spew"
)

type BridgeBasic struct {
	BridgeClient *eventbridge.Client
	BusName      string
}

// Create the client to connect with event Bridge
func EventManager(busName string) (BridgeBasic, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println(customErr.ErrEventBridgeClient.Error())
		return BridgeBasic{BridgeClient: nil, BusName: " "}, customErr.ErrEventBridgeClient
	}
	return BridgeBasic{BridgeClient: eventbridge.NewFromConfig(cfg), BusName: busName}, nil
}

// Send a new event
func (bridge BridgeBasic) SendEvent(source string, message interface{}) error {
	// Marshall the message to send into eventBridge
	msg, err := json.Marshal(message)
	if err != nil {
		return customErr.ErrMarsh
	}

	response, err := bridge.BridgeClient.PutEvents(context.TODO(),
		&eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					EventBusName: aws.String(bridge.BusName),
					Source:       aws.String(source),

					DetailType: aws.String("test"),
					Detail:     aws.String(string(msg)),
				},
			},
		},
	)
	if err != nil {
		log.Println("error sending the event to eventBridge" + err.Error())
		return customErr.ErrEventBrigePutEvent
	}
	spew.Dump(response.Entries)

	return nil
}
