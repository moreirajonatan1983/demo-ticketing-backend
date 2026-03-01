#!/bin/bash

export AWS_REGION="us-east-1"
export AWS_ACCESS_KEY_ID="test"
export AWS_SECRET_ACCESS_KEY="test"
export DYNAMODB_ENDPOINT="http://localhost:8000"
export LOCALSTACK_ENDPOINT="http://localhost:4566"
export AWS_ENDPOINT_URL="http://localhost:4566"

echo "Creating S3 Bucket 'ticketera-images-local'..."
aws --endpoint-url $LOCALSTACK_ENDPOINT s3 mb s3://ticketera-images-local || echo "Bucket already exists."

echo "Creating SQS DLQ 'ticket-purchased-dlq'..."
DLQ_URL=$(aws --endpoint-url $LOCALSTACK_ENDPOINT sqs create-queue \
    --queue-name ticket-purchased-dlq \
    --query 'QueueUrl' --output text 2>/dev/null) \
    || DLQ_URL=$(aws --endpoint-url $LOCALSTACK_ENDPOINT sqs get-queue-url --queue-name ticket-purchased-dlq --query 'QueueUrl' --output text)

DLQ_ARN=$(aws --endpoint-url $LOCALSTACK_ENDPOINT sqs get-queue-attributes \
    --queue-url $DLQ_URL --attribute-names QueueArn --query 'Attributes.QueueArn' --output text)

echo "Creating SQS queue 'ticket-purchased-queue' with RedrivePolicy..."
# Pasarlo correctamente escapando las comillas dobles que espera AWS
aws --endpoint-url $LOCALSTACK_ENDPOINT sqs create-queue \
    --queue-name ticket-purchased-queue \
    --attributes '{"RedrivePolicy":"{\"maxReceiveCount\":\"3\", \"deadLetterTargetArn\":\"'"$DLQ_ARN"'\"}"}' \
    || aws --endpoint-url $LOCALSTACK_ENDPOINT sqs set-queue-attributes \
    --queue-url $(aws --endpoint-url $LOCALSTACK_ENDPOINT sqs get-queue-url --queue-name ticket-purchased-queue --query 'QueueUrl' --output text) \
    --attributes '{"RedrivePolicy":"{\"maxReceiveCount\":\"3\", \"deadLetterTargetArn\":\"'"$DLQ_ARN"'\"}"}'

echo "Creating SNS topic 'ticketera-notifications'..."
aws --endpoint-url $LOCALSTACK_ENDPOINT sns create-topic \
    --name ticketera-notifications \
    || echo "SNS topic already exists."

echo "Creating EventBridge event bus 'ticketera-events'..."
aws --endpoint-url $LOCALSTACK_ENDPOINT events create-event-bus \
    --name ticketera-events \
    || echo "Event bus already exists."

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

echo "Creating ShowsTable..."
aws dynamodb create-table \
    --table-name ShowsTable \
    --attribute-definitions AttributeName=event_id,AttributeType=S \
    --key-schema AttributeName=event_id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url $DYNAMODB_ENDPOINT \
    || echo "Table ShowsTable already exists."

echo "Creating UsersTable (Auth)..."
aws dynamodb create-table \
    --table-name UsersTable \
    --attribute-definitions AttributeName=email,AttributeType=S \
    --key-schema AttributeName=email,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url $DYNAMODB_ENDPOINT \
    || echo "Table UsersTable already exists."


echo "Seeding..."
export EVENTS_TABLE_NAME="EventsTable"
export SHOWS_TABLE_NAME="ShowsTable"
export SEATS_TABLE_NAME="EventSeatsTable"
export TICKETS_TABLE_NAME="TicketsTable"
export AWS_ENDPOINT_URL="http://localhost:8000"

# Build and run the seed script!
go run seed_dynamo.go

echo "Done!"

echo "Setting up SAGA Step Functions..."
./setup_saga.sh
