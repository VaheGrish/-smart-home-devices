package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"smart-home-devices/internal/dynamoapi"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Device struct {
	ID         string `json:"id" dynamodbav:"id"`
	Mac        string `json:"mac" dynamodbav:"mac"`
	Name       string `json:"name" dynamodbav:"name"`
	Type       string `json:"type" dynamodbav:"type"`
	HomeID     string `json:"homeId" dynamodbav:"homeId"`
	CreatedAt  int64  `json:"createdAt" dynamodbav:"createdAt"`
	ModifiedAt int64  `json:"modifiedAt" dynamodbav:"modifiedAt"`
}

var (
	dynamoClient dynamoapi.DynamoAPI
	tableName    string
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
	tableName = os.Getenv("DEVICES_TABLE")
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received request: %s", req.Body)

	var device Device
	if err := json.Unmarshal([]byte(req.Body), &device); err != nil {
		log.Printf("Invalid JSON: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid JSON"}, nil
	}

	if device.ID == "" || device.Mac == "" || device.Name == "" || device.Type == "" {
		log.Printf("Missing required fields")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing required fields: id, mac, name, type must be present",
		}, nil
	}

	now := time.Now().UnixMilli()
	device.CreatedAt = now
	device.ModifiedAt = now

	item, err := attributevalue.MarshalMap(device)
	if err != nil {
		log.Printf("Error marshaling item: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Marshal error"}, nil
	}

	log.Printf("Putting item into DynamoDB table %s: %+v", tableName, device)

	_, err = dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		log.Printf("DynamoDB PutItem error: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "DynamoDB error"}, nil
	}

	log.Printf("Device created successfully: %s", device.ID)
	return events.APIGatewayProxyResponse{StatusCode: 201, Body: "Device created"}, nil
}

func main() {
	lambda.Start(handler)
}
