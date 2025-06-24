package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"smart-home-devices/internal/dynamoapi"
	"smart-home-devices/internal/service"
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
	if tableName == "" {
		log.Fatal("DEVICES_TABLE environment variable not set")
	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]
	if id == "" {
		log.Println("Missing device ID in request")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing device ID",
		}, nil
	}

	err := service.DeleteDevice(ctx, dynamoClient, tableName, id)
	if err != nil {
		log.Printf("Failed to delete device with id %s: %v", id, err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to delete device",
		}, nil
	}

	log.Printf("Device deleted successfully: %s", id)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Device deleted",
	}, nil
}

func main() {
	lambda.Start(handler)
}
