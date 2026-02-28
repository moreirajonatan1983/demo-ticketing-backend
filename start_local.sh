#!/bin/bash

# Make sure Docker is running
echo "Starting DynamoDB Local Container..."
docker-compose up -d

echo "Running migrations / seeds..."
sleep 2 # wait for Dynamodb
./setup_local_dynamo.sh

echo "Building all Lambdas..."
cd lambdas/events-lambda && GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go
cd ../seats-lambda && GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go
cd ../checkout-lambda && GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go
cd ../tickets-lambda && GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go
cd ../..

echo "Start complete!"
echo "Run each Lambda API using SAM local in different terminal tabs:"
echo "--- Tab 1: cd lambdas/events-lambda && sam local start-api -p 3000 --env-vars ../../env.json"
echo "--- Tab 2: cd lambdas/checkout-lambda && sam local start-api -p 3004 --env-vars ../../env.json"
echo "--- Tab 3: cd lambdas/seats-lambda && sam local start-api -p 3005 --env-vars ../../env.json"
echo "--- Tab 4: cd lambdas/tickets-lambda && sam local start-api -p 3006 --env-vars ../../env.json"
