package dynamoapi

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type MockDynamoClient struct {
	PutItemFunc    func(ctx context.Context, input *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItemFunc    func(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	UpdateItemFunc func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	DeleteItemFunc func(ctx context.Context, input *dynamodb.DeleteItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

func (m *MockDynamoClient) PutItem(ctx context.Context, input *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if m.PutItemFunc != nil {
		return m.PutItemFunc(ctx, input, opts...)
	}
	return &dynamodb.PutItemOutput{}, nil
}

func (m *MockDynamoClient) GetItem(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if m.GetItemFunc != nil {
		return m.GetItemFunc(ctx, input, opts...)
	}
	return &dynamodb.GetItemOutput{}, nil
}

func (m *MockDynamoClient) UpdateItem(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	if m.UpdateItemFunc != nil {
		return m.UpdateItemFunc(ctx, input, opts...)
	}
	return &dynamodb.UpdateItemOutput{}, nil
}

func (m *MockDynamoClient) DeleteItem(ctx context.Context, input *dynamodb.DeleteItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	if m.DeleteItemFunc != nil {
		return m.DeleteItemFunc(ctx, input, opts...)
	}
	return &dynamodb.DeleteItemOutput{}, nil
}
