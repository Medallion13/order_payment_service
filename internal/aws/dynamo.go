package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

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

func (table TableBasics) UpdateInfo(tableName string, keyName string, upInfo interface{}) error {

	//! old
	// Obtein the key value
	// key := infoMap[keyName]
	// fmt.Println(infoMap[keyName])

	// var some expression.updateder

	// // Contrut de update builder to send the update information
	// for attrName := range infoMap {
	// 	if attrName != keyName && attrName != "CreateAt" {
	// 		some.Set(expression.Name(attrName), expression.Value(infoMap[attrName]))
	// 		a, b := infoMap[attrName]
	// 		fmt.Println(attrName)
	// 		fmt.Println(infoMap[attrName])
	// 		fmt.Println(a)
	// 		fmt.Println(b)
	// 		fmt.Printf("Type of value: %T\n", infoMap[attrName])
	// 	} else {
	// 		fmt.Printf("Skipping update for attribute: %s\n", attrName)
	// 	}
	// }
	// some.Set(expression.Name("UpdateAt"), expression.Value(time.Now().Format(time.RFC822)))

	// fmt.Printf("Type of some: %T ", some)

	// updateBuilder = updateBuilder.WithUpdate(some)
	// updateBuilder = updateBuilder.WithUpdate(
	// )

	// Building the updete objet to the request to Dynamo
	// update, err := expression.NewBuilder().WithUpdate(some).Build()
	// if err != nil {
	// 	fmt.Println(fmt.Sprintf("Error: %s", err))
	// 	return customErr.ErrBuildingExpression
	// }

	// updateInput := dynamodb.UpdateItemInput{
	// 	TableName: aws.String(tableName),
	// 	Key: map[string]types.AttributeValue{
	// 		keyName: key,
	// 	},
	// 	UpdateExpression:          update.Update(),
	// 	ExpressionAttributeNames:  update.Names(),
	// 	ExpressionAttributeValues: update.Values(),
	// }

	//! actual
	// convert the struct into a map to have access to all values
	upInfoMap, err := StructToMap(upInfo)
	if err != nil {
		return err
	}

	// Obtain the key value
	keyvalue := upInfoMap[keyName]

	// create a expression builder object to construct the update object
	update := expression.Set(expression.Name("UpdateAt"), expression.Value(time.Now().Format(time.RFC822)))

	for mapKey, mapValue := range upInfoMap {

		if mapKey != keyName && mapKey != "CreateAt" && !utils.IsEmpty(mapValue) {
			update.Set(expression.Name(mapKey), expression.Value(fmt.Sprint(mapValue)))
		}
	}

	fmt.Println(update)

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return customErr.ErrBuildingExpression
	}

	fmt.Printf("Length: %v\n", len(expr.Names()))
	for key, value := range expr.Names() {
		fmt.Printf("%s: %v\n", key, value)
	}

	updateInput := dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			keyName: &types.AttributeValueMemberS{Value: fmt.Sprint(keyvalue)},
		},
		UpdateExpression:         expr.Update(),
		ExpressionAttributeNames: expr.Names(),
	}

	return nil
}

func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}
