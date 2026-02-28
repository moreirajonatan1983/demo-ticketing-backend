#!/bin/bash

# Make sure Docker is running
echo "Starting DynamoDB Local Container..."
docker-compose up -d

echo "Running migrations / seeds..."
sleep 2 # wait for Dynamodb
./setup_local_dynamo.sh

echo "Building all Lambdas..."
cd lambdas/events-aws-lambda && GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go
cd ../seats-aws-lambda && GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go
cd ../checkout-aws-lambda && GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go
cd ../tickets-aws-lambda && GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go
cd ../..

echo "Start complete!"
echo "Run each Lambda API using SAM local in different terminal tabs:"
echo "--- Tab 1: cd lambdas/events-aws-lambda && sam local start-api -p 3000 --env-vars ../../env.json"
echo "--- Tab 2: cd lambdas/checkout-aws-lambda && sam local start-api -p 3004 --env-vars ../../env.json"
echo "--- Tab 3: cd lambdas/seats-aws-lambda && sam local start-api -p 3005 --env-vars ../../env.json"
echo "--- Tab 4: cd lambdas/tickets-aws-lambda && sam local start-api -p 3006 --env-vars ../../env.json"
