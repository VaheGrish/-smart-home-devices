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
- Serverless Framework installed (`npm install -g serverless`)
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
sls invoke -f create --data '{
  "body": "{\"id\":\"1\",\"mac\":\"AA:BB:CC:DD:EE:FF\",\"name\":\"Thermostat\",\"type\":\"thermostat\",\"homeId\":\"home-123\"}"
}'

