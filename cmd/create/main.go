package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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
	dynamoClient *dynamodb.Client
	tableName    string
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Unable to load AWS config: " + err.Error())
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
	tableName = os.Getenv("DEVICES_TABLE")
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var device Device

	if err := json.Unmarshal([]byte(req.Body), &device); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid JSON"}, nil
	}

	if device.ID == "" || device.Mac == "" || device.Name == "" || device.Type == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Missing required fields"}, nil
	}

	now := time.Now().UnixMilli()
	device.CreatedAt = now
	device.ModifiedAt = now

	item, err := attributevalue.MarshalMap(device)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Marshaling error"}, nil
	}

	_, err = dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("DynamoDB error: %v", err)}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 201, Body: "Device created"}, nil
}

func main() {
	lambda.Start(handler)
}
