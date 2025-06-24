package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"smart-home-devices/internal/dynamoapi"
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
	ErrMissingID      = errors.New("missing device ID")
	ErrDeviceNotFound = errors.New("device not found")
	ErrMissingFields  = errors.New("missing required fields: id, mac, name, type must be present")
)

func CreateDevice(ctx context.Context, client dynamoapi.DynamoAPI, table string, device *Device) error {
	if device.ID == "" || device.Mac == "" || device.Name == "" || device.Type == "" {
		return ErrMissingFields
	}

	now := time.Now().UnixMilli()
	device.CreatedAt = now
	device.ModifiedAt = now

	item, err := attributevalue.MarshalMap(device)
	if err != nil {
		return err
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item:      item,
	})
	return err
}

func GetDeviceByID(ctx context.Context, client dynamoapi.DynamoAPI, tableName, id string) (string, error) {
	if id == "" {
		return "", ErrMissingID
	}

	resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return "", fmt.Errorf("dynamodb GetItem error: %w", err)
	}

	if resp.Item == nil || len(resp.Item) == 0 {
		return "", ErrDeviceNotFound
	}

	var device Device
	if err := attributevalue.UnmarshalMap(resp.Item, &device); err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	data, err := json.Marshal(device)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	return string(data), nil
}

func UpdateDevice(ctx context.Context, client dynamoapi.DynamoAPI, tableName string, updateReq map[string]interface{}) error {
	id, ok := updateReq["id"].(string)
	if !ok || id == "" {
		return ErrMissingID
	}

	expr := "SET modifiedAt = :mod"
	values := map[string]types.AttributeValue{
		":mod": &types.AttributeValueMemberN{Value: fmt.Sprint(time.Now().UnixMilli())},
	}
	names := map[string]string{}

	if name, ok := updateReq["name"].(string); ok {
		expr += ", #N = :name"
		values[":name"] = &types.AttributeValueMemberS{Value: name}
		names["#N"] = "name"
	}
	if mac, ok := updateReq["mac"].(string); ok {
		expr += ", mac = :mac"
		values[":mac"] = &types.AttributeValueMemberS{Value: mac}
	}
	if typ, ok := updateReq["type"].(string); ok {
		expr += ", #T = :type"
		values[":type"] = &types.AttributeValueMemberS{Value: typ}
		names["#T"] = "type"
	}
	if homeId, ok := updateReq["homeId"].(string); ok {
		expr += ", homeId = :homeId"
		values[":homeId"] = &types.AttributeValueMemberS{Value: homeId}
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: id}},
		UpdateExpression:          aws.String(expr),
		ExpressionAttributeValues: values,
	}
	if len(names) > 0 {
		input.ExpressionAttributeNames = names
	}

	_, err := client.UpdateItem(ctx, input)
	return err
}

func DeleteDevice(ctx context.Context, client dynamoapi.DynamoAPI, tableName, id string) error {
	if id == "" {
		return ErrMissingID
	}
	_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}
