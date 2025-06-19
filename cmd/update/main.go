package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	dynamoClient *dynamodb.Client
	tableName    string
)

type UpdateDeviceRequest struct {
	ID     string  `json:"id"`
	Mac    *string `json:"mac,omitempty"`
	Name   *string `json:"name,omitempty"`
	Type   *string `json:"type,omitempty"`
	HomeID *string `json:"homeId,omitempty"`
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Unable to load AWS config: " + err.Error())
	}
	dynamoClient = dynamodb.NewFromConfig(cfg)
	tableName = os.Getenv("DEVICES_TABLE")
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var updateReq UpdateDeviceRequest
	err := json.Unmarshal([]byte(req.Body), &updateReq)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid JSON"}, nil
	}

	if updateReq.ID == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Missing device ID"}, nil
	}

	updateExpr := "SET modifiedAt = :mod"
	exprAttrValues := map[string]types.AttributeValue{
		":mod": &types.AttributeValueMemberN{Value: strconv.FormatInt(time.Now().UnixMilli(), 10)},
	}
	exprAttrNames := map[string]string{}

	if updateReq.Name != nil {
		updateExpr += ", #N = :name"
		exprAttrValues[":name"] = &types.AttributeValueMemberS{Value: *updateReq.Name}
		exprAttrNames["#N"] = "name"
	}
	if updateReq.Mac != nil {
		updateExpr += ", mac = :mac"
		exprAttrValues[":mac"] = &types.AttributeValueMemberS{Value: *updateReq.Mac}
	}
	if updateReq.Type != nil {
		updateExpr += ", #T = :type"
		exprAttrValues[":type"] = &types.AttributeValueMemberS{Value: *updateReq.Type}
		exprAttrNames["#T"] = "type" 
	}
	if updateReq.HomeID != nil {
		updateExpr += ", homeId = :homeId"
		exprAttrValues[":homeId"] = &types.AttributeValueMemberS{Value: *updateReq.HomeID}
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: updateReq.ID}},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeValues: exprAttrValues,
	}

	if len(exprAttrNames) > 0 {
		input.ExpressionAttributeNames = exprAttrNames
	}

	_, err = dynamoClient.UpdateItem(ctx, input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Update error: %v", err)}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Device updated"}, nil
}

func main() {
	lambda.Start(handler)
}
