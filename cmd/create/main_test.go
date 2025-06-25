package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"

	"smart-home-devices/internal/dynamoapi"
	"smart-home-devices/internal/service"
)

func HandlerWithClient(client dynamoapi.DynamoAPI, tableName string) func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	svc := service.NewDeviceService(client, tableName)
	return HandlerWithService(svc)
}

func HandlerWithService(svc service.DeviceService) func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var device service.Device
		if err := json.Unmarshal([]byte(req.Body), &device); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid JSON"}, nil
		}

		err := svc.Create(ctx, &device)
		if err != nil {
			if err == service.ErrMissingFields {
				return events.APIGatewayProxyResponse{StatusCode: 400, Body: err.Error()}, nil
			}
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, nil
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

	handler := HandlerWithClient(mockClient, "DevicesTable")

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
	handler := HandlerWithClient(mockClient, "DevicesTable")

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
	handler := HandlerWithClient(mockClient, "DevicesTable")

	request := events.APIGatewayProxyRequest{
		Body: `{"id":"","mac":"AA:BB:CC:DD:EE:FF"}`,
	}

	resp, err := handler(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "missing required fields: id, mac, name, type must be present", resp.Body)
}
