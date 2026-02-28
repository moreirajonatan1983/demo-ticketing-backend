#!/bin/bash

# Make sure Docker is running
echo "Starting DynamoDB Local Container..."
cd .. && docker-compose up -d && cd scripts

echo "Running migrations / seeds..."
sleep 2 # wait for Dynamodb
./setup_local_dynamo.sh

echo "Building all Lambdas..."
cd ../lambdas/events-aws-lambda && GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go
cd ../media-aws-lambda && GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go
cd ../shows-aws-lambda && GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go
cd ../seats-aws-lambda && GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go
cd ../checkout-aws-lambda && GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go
cd ../tickets-aws-lambda && GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/api/main.go
cd ../..

echo "Start complete!"
echo "Run each Lambda API using SAM local in different terminal tabs:"
echo "--- Tab 1: cd lambdas/events-aws-lambda && sam local start-api -p 3000 --env-vars ../../env.json"
echo "--- Tab 2: cd lambdas/shows-aws-lambda && sam local start-api -p 3007 --env-vars ../../env.json"
echo "--- Tab 3: cd lambdas/checkout-aws-lambda && sam local start-api -p 3004 --env-vars ../../env.json"
echo "--- Tab 4: cd lambdas/seats-aws-lambda && sam local start-api -p 3005 --env-vars ../../env.json"
echo "--- Tab 5: cd lambdas/tickets-aws-lambda && sam local start-api -p 3006 --env-vars ../../env.json"
echo "--- Tab 6: cd lambdas/media-aws-lambda && sam local start-api -p 3008 --env-vars ../../env.json"
