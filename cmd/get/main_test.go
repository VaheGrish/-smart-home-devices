package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"

	"smart-home-devices/internal/dynamoapi"
)

func HandlerWithClient(client dynamoapi.DynamoAPI) func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		id := req.PathParameters["id"]
		if id == "" {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Missing device ID"}, nil
		}

		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String("DevicesTable"),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: id},
			},
		})
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "DynamoDB error: " + err.Error()}, nil
		}

		if resp.Item == nil || len(resp.Item) == 0 {
			return events.APIGatewayProxyResponse{StatusCode: 404, Body: "Device not found"}, nil
		}

		var device Device
		if err := attributevalue.UnmarshalMap(resp.Item, &device); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Unmarshal error"}, nil
		}

		body, err := json.Marshal(device)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Marshal error"}, nil
		}

		return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(body)}, nil
	}
}

func TestGetHandler_Success(t *testing.T) {
	device := Device{
		ID:         "1",
		Mac:        "AA:BB:CC:DD:EE:FF",
		Name:       "Thermostat",
		Type:       "thermostat",
		HomeID:     "home-123",
		CreatedAt:  1234567890,
		ModifiedAt: 1234567890,
	}

	item, _ := attributevalue.MarshalMap(device)

	mockClient := &dynamoapi.MockDynamoClient{
		GetItemFunc: func(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{Item: item}, nil
		},
	}

	handler := HandlerWithClient(mockClient)
	req := events.APIGatewayProxyRequest{PathParameters: map[string]string{"id": "1"}}

	resp, err := handler(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var got Device
	_ = json.Unmarshal([]byte(resp.Body), &got)
	assert.Equal(t, device.ID, got.ID)
	assert.Equal(t, device.Name, got.Name)
}

func TestGetHandler_NotFound(t *testing.T) {
	mockClient := &dynamoapi.MockDynamoClient{
		GetItemFunc: func(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{Item: nil}, nil
		},
	}

	handler := HandlerWithClient(mockClient)
	req := events.APIGatewayProxyRequest{PathParameters: map[string]string{"id": "non-existent"}}

	resp, err := handler(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
	assert.Equal(t, "Device not found", resp.Body)
}

func TestGetHandler_MissingID(t *testing.T) {
	mockClient := &dynamoapi.MockDynamoClient{}

	handler := HandlerWithClient(mockClient)
	req := events.APIGatewayProxyRequest{PathParameters: map[string]string{}}

	resp, err := handler(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "Missing device ID", resp.Body)
}
