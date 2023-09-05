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
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	customErr "github.com/NicoCodes13/order_payment_service/internal/errors"
	"github.com/NicoCodes13/order_payment_service/internal/utils"
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
func (table TableBasics) PutItem(info interface{}) error {
	// transform the info into map[string]interface{}
	item, err := attributevalue.MarshalMap(info)
	if err != nil {
		return customErr.ErrMarsh
	}

	// Send the information to create a new item in the table
	_, err = table.DynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(table.TableName),
		Item:      item,
	})
	if err != nil {
		log.Printf("couldn't add item to table. Here's why: %v\n", err)
		return err
	}
	return err
}

// Update a exist item in the dynamo table
func (table TableBasics) UpdateInfo(keyName string, upInfo interface{}) error {
	// convert the struct into a map to have access to all values
	upInfoMap, err := utils.StructToMap(upInfo)
	if err != nil {
		return err
	}

	// Obtain the key value
	keyvalue := upInfoMap[keyName]

	// create a expression builder object to construct the update object
	update := expression.Set(expression.Name("UpdateAt"), expression.Value(time.Now().Format(time.RFC822)))

	// it goes through the object and stores the names and values
	// leaving out the empty ones and the key
	for mapKey, mapValue := range upInfoMap {

		if mapKey != keyName && mapKey != "CreateAt" && !utils.IsEmpty(mapValue) {
			update.Set(expression.Name(mapKey), expression.Value(mapValue))
		}
	}
	// Condition to prevent creation of new records
	condition := expression.Name(keyName).Equal(expression.Value(&types.AttributeValueMemberS{
		Value: fmt.Sprint(keyvalue),
	}))

	// Building the expretion to construct te object to update
	expr, err := expression.NewBuilder().WithUpdate(update).WithCondition(condition).Build()
	if err != nil {
		log.Println(err)
		return customErr.ErrBuildingExpression
	}

	// Object to make the input process
	updateInput := dynamodb.UpdateItemInput{
		TableName: aws.String(table.TableName),
		Key: map[string]types.AttributeValue{
			keyName: &types.AttributeValueMemberS{Value: fmt.Sprint(keyvalue)},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	}

	// Making an update of the information
	out, err := table.DynamoClient.UpdateItem(context.TODO(), &updateInput)
	if err != nil {
		log.Printf("Couldn't update\n %v", err)
		return customErr.ErrUpdateDynamo
	}

	// Print the response's attributes
	if len(out.Attributes) > 0 {
		log.Println("Updated Values: ")
		for attributeName, attributeValue := range out.Attributes {
			log.Printf("%s: %v\n", attributeName, attributeValue)
		}
	} else {
		log.Println("No attributes returned in the response.")
	}

	return nil
}

// Get item
func (table TableBasics) GetItem(keyName string, keyValue string, out interface{}) error {
	response, err := table.DynamoClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &table.TableName,
		Key: map[string]types.AttributeValue{
			keyName: &types.AttributeValueMemberS{Value: keyValue},
		},
	})
	if err != nil {
		return err
	}

	err = attributevalue.UnmarshalMap(response.Item, &out)
	if err != nil {
		return err
	}
	return nil
}
