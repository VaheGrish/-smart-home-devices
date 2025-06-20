# Smart Home Devices Service

This is a serverless AWS Lambda-based service written in Go for managing smart home devices. It supports basic CRUD operations and integrates with DynamoDB and SQS.

## Features

- Create, Get, Update, Delete devices
- DynamoDB used as a data store
- SQS queue to handle asynchronous updates
- Built with AWS Lambda using Go runtime (provided.al2)
- Serverless Framework for deployment automation

## Prerequisites

- Go 1.20+ installed
- AWS CLI configured with appropriate permissions
- Serverless Framework v3 installed (`npm install -g serverless@3`)
- AWS account with DynamoDB and SQS permissions

## Setup

1. Clone the repository:

```bash
git clone https://github.com/VaheGrish/-smart-home-devices.git
cd smart-home-devices

make build         # Build all Lambda functions
make deploy        # Build & deploy all functions
make deploy-<fn>   # Build & deploy a specific function 
make clean         # Clean up build artifacts

## Invoke a function with sample input:

## create
sls invoke -f create --data '{
  "body": "{\"id\":\"3\",\"mac\":\"AA:BB:CC:DD:EE:FF\",\"name\":\"Thermostat\",\"type\":\"thermostat\",\"homeId\":\"home-123\"}"
}'

## get
sls invoke -f get --data '{"pathParameters": {"id": "1"}}'

## update
sls invoke -f update -d '{
  "body": "{\"id\":\"1\", \"name\":\"Updated Device\", \"type\":\"sensor\"}"
}'

## delete
sls invoke -f delete -d '{
  "pathParameters": { "id": "2" }
}'

## sqs
aws sqs send-message   --queue-url https://sqs.<region>.amazonaws.com/<account-id>/DeviceAssociationQueue   --message-body '{"id":"1", "homeId":"home-NEW-1"}'

# Run all unit tests recursively with verbose output
go test ./... -v

# Run tests only in the 'cmd/create' package
go test ./cmd/create -v
