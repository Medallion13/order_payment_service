package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	customErr "github.com/NicoCodes13/order_payment_service/internal/errors"
)

type TableBasics struct {
	DynamoClient *dynamodb.Client
	TableName    string
}

func DynamoClient(tableName string) (TableBasics, error) {
	config, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Print(customErr.ErrAPIClient.Error())
		err = customErr.ErrAPIClient
		return TableBasics{DynamoClient: nil, TableName: ""}, err
	}
	return TableBasics{DynamoClient: dynamodb.NewFromConfig(config), TableName: tableName}, nil
}

// Check if the table exists
func (basics TableBasics) TableExists() (bool, error) {
	exist := true
	_, err := basics.DynamoClient.DescribeTable(
		context.TODO(), &dynamodb.DescribeTableInput{TableName: aws.String(basics.TableName)},
	)
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			log.Printf("Table %v does not exist. \n", basics.TableName)
			err = nil
		} else {
			log.Printf("Couldn't determine existence of table %v, Here's why: %v\n", basics.TableName, err)
		}
		exist = false
	}
	return exist, err
}

// Creating a new item in the table
func (table TableBasics) PutInfo(info interface{}) error {
	// transform the info into map[string]interface{}
	item, err := attributevalue.MarshalMap(info)
	if err != nil {
		return customErr.ErrMarsh
	}
	_, err = table.DynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(table.TableName),
		Item:      item,
	})
	if err != nil {
		log.Printf("couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

func (table TableBasics) UpdateInfo(tableName string, keyName string, info interface{}) error {
	// transform the info into map[string]interface{}
	infoMap, err := attributevalue.MarshalMap(info)
	if err != nil {
		return customErr.ErrMarsh
	}

	// Contrut de update builder to send the update information
	updateBuilder := expression.NewBuilder()
	for attrName := range infoMap {
		updateBuilder = updateBuilder.WithUpdate(
			expression.Set(expression.Name(attrName), expression.Value(infoMap[attrName])),
		)
	}

	updateBuilder = updateBuilder.WithUpdate(
		expression.Set(expression.Name("UpdateAt"), expression.Value(time.Now().Format(time.RFC822))),
	)

	// Building the updete objet to the request to Dynamo
	update, err := updateBuilder.Build()
	if err != nil {
		return customErr.ErrBuildingExpression
	}

	fmt.Print(update)

	// Define the key condition

	// // Define the update input parameters
	// updateInput := &dynamodb.UpdateItemInput{
	// 	TableName: aws.String(tableName),
	// 	Key: map,
	// }

	// var response *dynamodb.UpdateItemOutput
	// var attributes map[string]map[string]interface{}
	return nil
}
