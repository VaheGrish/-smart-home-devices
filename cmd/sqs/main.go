package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"smart-home-devices/internal/dynamoapi"
)

var (
	dynamoClient dynamoapi.DynamoAPI
	tableName    string
)

type UpdateMessage struct {
	ID     string `json:"id"`
	HomeID string `json:"homeId"`
}

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

func handler(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		log.Printf("Received message body: %s", record.Body)

		var msg UpdateMessage
		if err := json.Unmarshal([]byte(record.Body), &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		now := time.Now().UnixMilli()
		_, err := dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: msg.ID},
			},
			UpdateExpression: aws.String("SET homeId = :h, modifiedAt = :m"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":h": &types.AttributeValueMemberS{Value: msg.HomeID},
				":m": &types.AttributeValueMemberN{Value: fmt.Sprint(now)},
			},
		})
		if err != nil {
			log.Printf("Failed to update device %s: %v", msg.ID, err)
			continue
		}
		log.Printf("Updated device: %s", msg.ID)
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
