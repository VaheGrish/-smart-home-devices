package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"smart-home-devices/internal/dynamoapi"
)

type Device struct {
	ID         string `json:"id"`
	Mac        string `json:"mac"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	HomeID     string `json:"homeId"`
	CreatedAt  int64  `json:"createdAt"`
	ModifiedAt int64  `json:"modifiedAt"`
}

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
	if id == "" {
		log.Println("Missing device ID in request")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing device ID",
		}, nil
	}

	resp, err := dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		log.Printf("DynamoDB GetItem error for id %s: %v", id, err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	}

	if resp.Item == nil || len(resp.Item) == 0 {
		log.Printf("Device not found: %s", id)
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Device not found",
		}, nil
	}

	var device Device
	if err := attributevalue.UnmarshalMap(resp.Item, &device); err != nil {
		log.Printf("Unmarshal error for id %s: %v", id, err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	}

	body, err := json.Marshal(device)
	if err != nil {
		log.Printf("Marshal error for id %s: %v", id, err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}
