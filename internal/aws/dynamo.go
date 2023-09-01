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
		fmt.Println("soy el error 1")
		return customErr.ErrMarsh
	}
	// Obtein the key value
	key := infoMap[keyName]
	fmt.Println(infoMap[keyName])

	var some expression.UpdateBuilder

	// Contrut de update builder to send the update information
	for attrName := range infoMap {
		if attrName != keyName && attrName != "CreateAt" {
			some.Set(expression.Name(attrName), expression.Value(infoMap[attrName]))
			a, b := infoMap[attrName]
			fmt.Println(attrName)
			fmt.Println(infoMap[attrName])
			fmt.Println(a)
			fmt.Println(b)
			fmt.Printf("Type of value: %T\n", infoMap[attrName])
		} else {
			fmt.Printf("Skipping update for attribute: %s\n", attrName)
		}
	}
	some.Set(expression.Name("UpdateAt"), expression.Value(time.Now().Format(time.RFC822)))

	fmt.Printf("Type of some: %T ", some)

	// updateBuilder = updateBuilder.WithUpdate(some)
	// updateBuilder = updateBuilder.WithUpdate(
	// )

	// Building the updete objet to the request to Dynamo
	update, err := expression.NewBuilder().WithUpdate(some).Build()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: %s", err))
		return customErr.ErrBuildingExpression
	}

	// updateInput := dynamodb.UpdateItemInput{
	// 	TableName: aws.String(tableName),
	// 	Key: map[string]types.AttributeValue{
	// 		keyName: key,
	// 	},
	// 	UpdateExpression:          update.Update(),
	// 	ExpressionAttributeNames:  update.Names(),
	// 	ExpressionAttributeValues: update.Values(),
	// }

	fmt.Printf("Length: %v\n", len(update.Names()))
	for key, value := range update.Names() {
		fmt.Printf("%s: %v\n", key, value)
	}
	fmt.Println(key)

	return nil
}
