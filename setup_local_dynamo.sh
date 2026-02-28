#!/bin/bash

export AWS_REGION="us-east-1"
export AWS_ACCESS_KEY_ID="test"
export AWS_SECRET_ACCESS_KEY="test"
export DYNAMODB_ENDPOINT="http://localhost:8000"

echo "Creating EventsTable..."
aws dynamodb create-table \
    --table-name EventsTable \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url $DYNAMODB_ENDPOINT \
    || echo "Table EventsTable already exists."

echo "Creating EventSeatsTable..."
aws dynamodb create-table \
    --table-name EventSeatsTable \
    --attribute-definitions AttributeName=event_id,AttributeType=S AttributeName=seat_id,AttributeType=S \
    --key-schema AttributeName=event_id,KeyType=HASH AttributeName=seat_id,KeyType=RANGE \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url $DYNAMODB_ENDPOINT \
    || echo "Table EventSeatsTable already exists."

echo "Creating TicketsTable..."
aws dynamodb create-table \
    --table-name TicketsTable \
    --attribute-definitions AttributeName=user_id,AttributeType=S AttributeName=ticket_id,AttributeType=S \
    --key-schema AttributeName=user_id,KeyType=HASH AttributeName=ticket_id,KeyType=RANGE \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url $DYNAMODB_ENDPOINT \
    || echo "Table TicketsTable already exists."

echo "Seeding..."
export EVENTS_TABLE_NAME="EventsTable"
export SEATS_TABLE_NAME="EventSeatsTable"
export TICKETS_TABLE_NAME="TicketsTable"
export AWS_ENDPOINT_URL="http://localhost:8000"

# Build and run the seed script!
cd scripts
go run seed_dynamo.go
cd ..

echo "Done!"
