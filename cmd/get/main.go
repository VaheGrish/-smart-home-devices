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
		log.Fatalf("unable to load AWS config: %v", err)
	}
	dynamoClient = dynamodb.NewFromConfig(cfg)

	tableName = os.Getenv("DEVICES_TABLE")
	if tableName == "" {
		log.Fatalf("DEVICES_TABLE environment variable is not set")
	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]
	log.Printf("Received request for device ID: %s", id)

	data, err := service.GetDeviceByID(ctx, dynamoClient, tableName, id)
	if err != nil {
		log.Printf("Error getting device %s: %v", id, err)
		switch err {
		case service.ErrMissingID:
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: err.Error()}, nil
		case service.ErrDeviceNotFound:
			return events.APIGatewayProxyResponse{StatusCode: 404, Body: err.Error()}, nil
		default:
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal server error"}, nil
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       data,
	}, nil
}

func main() {
	lambda.Start(handler)
}
