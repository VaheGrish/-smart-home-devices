package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"

	"smart-home-devices/internal/dynamoapi"
)

func HandlerWithClient(client dynamoapi.DynamoAPI) func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String("DevicesTable"),
			Item:      item,
		})
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "DynamoDB error: " + err.Error()}, nil
		}

		return events.APIGatewayProxyResponse{StatusCode: 201, Body: "Device created"}, nil
	}
}

func TestCreateHandler_Success(t *testing.T) {
	mockClient := &dynamoapi.MockDynamoClient{
		PutItemFunc: func(ctx context.Context, input *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	}

	handler := HandlerWithClient(mockClient)

	request := events.APIGatewayProxyRequest{
		Body: `{"id":"1","mac":"AA:BB:CC:DD:EE:FF","name":"Thermostat","type":"thermostat","homeId":"home-123"}`,
	}

	resp, err := handler(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	assert.Equal(t, "Device created", resp.Body)
}

func TestCreateHandler_InvalidJSON(t *testing.T) {
	mockClient := &dynamoapi.MockDynamoClient{}
	handler := HandlerWithClient(mockClient)

	request := events.APIGatewayProxyRequest{
		Body: `invalid json`,
	}

	resp, err := handler(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "Invalid JSON", resp.Body)
}

func TestCreateHandler_MissingFields(t *testing.T) {
	mockClient := &dynamoapi.MockDynamoClient{}
	handler := HandlerWithClient(mockClient)

	request := events.APIGatewayProxyRequest{
		Body: `{"id":"","mac":"AA:BB:CC:DD:EE:FF"}`,
	}

	resp, err := handler(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "Missing required fields", resp.Body)
}
