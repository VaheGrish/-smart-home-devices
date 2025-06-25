package main

import (
	"context"
	"encoding/json"
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

		updateData := map[string]interface{}{
			"id":     msg.ID,
			"homeId": msg.HomeID,
		}

		svc := service.NewDeviceService(dynamoClient, tableName)
		err := svc.Update(ctx, updateData)

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
