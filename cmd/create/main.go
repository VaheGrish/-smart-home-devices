package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"smart-home-devices/internal/dynamoapi"
	"smart-home-devices/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

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

	var device service.Device
	if err := json.Unmarshal([]byte(req.Body), &device); err != nil {
		log.Printf("Invalid JSON: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid JSON"}, nil
	}

	err := service.CreateDevice(ctx, dynamoClient, tableName, &device)
	if err != nil {
		if err == service.ErrMissingFields {
			log.Printf("Validation error: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: err.Error()}, nil
		}
		log.Printf("DynamoDB error: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, nil
	}

	log.Printf("Device created successfully: %s", device.ID)
	return events.APIGatewayProxyResponse{StatusCode: 201, Body: "Device created"}, nil
}

func main() {
	lambda.Start(handler)
}
