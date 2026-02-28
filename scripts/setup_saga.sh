#!/bin/bash
# setup_saga.sh - Registers SAGA Lambda functions and Step Functions state machine in LocalStack

export AWS_REGION="us-east-1"
export AWS_ACCESS_KEY_ID="test"
export AWS_SECRET_ACCESS_KEY="test"
export LOCALSTACK="http://localhost:4566"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/.."
LAMBDAS_DIR="$BACKEND_DIR/lambdas"
SAGA_DIR="$BACKEND_DIR/saga"

echo ">>> [1/3] Building SAGA Lambda binaries..."

echo "  Building seats-aws-lambda..."
cd "$LAMBDAS_DIR/seats-aws-lambda"
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap ./cmd/api
zip -j seats.zip bootstrap

echo "  Building checkout-aws-lambda..."
cd "$LAMBDAS_DIR/checkout-aws-lambda"
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap ./cmd/api
zip -j checkout.zip bootstrap

echo "  Building tickets-aws-lambda..."
cd "$LAMBDAS_DIR/tickets-aws-lambda"
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap ./cmd/api
zip -j tickets.zip bootstrap

echo ""
echo ">>> [2/3] Registering Lambda functions in LocalStack..."

LAMBDA_NAMES=("saga-reserve-seat" "saga-release-seat" "saga-process-payment" "saga-create-ticket")
for name in "${LAMBDA_NAMES[@]}"; do
  STATUS=$(aws --endpoint-url $LOCALSTACK lambda get-function --function-name "$name" 2>&1)
  if echo "$STATUS" | grep -q "ResourceNotFoundException"; then
    echo "  Creating Lambda: $name"
    case "$name" in
      "saga-reserve-seat")
        aws --endpoint-url $LOCALSTACK lambda create-function \
          --function-name "$name" --runtime provided.al2023 --architectures arm64 \
          --handler bootstrap --role arn:aws:iam::000000000000:role/dummy \
          --zip-file "fileb://$LAMBDAS_DIR/seats-aws-lambda/seats.zip" \
          --environment "Variables={SEATS_TABLE_NAME=EventSeatsTable,AWS_ENDPOINT_URL=http://host.docker.internal:8000}" > /dev/null
        ;;
      "saga-release-seat")
        aws --endpoint-url $LOCALSTACK lambda create-function \
          --function-name "$name" --runtime provided.al2023 --architectures arm64 \
          --handler bootstrap --role arn:aws:iam::000000000000:role/dummy \
          --zip-file "fileb://$LAMBDAS_DIR/seats-aws-lambda/seats.zip" \
          --environment "Variables={SEATS_TABLE_NAME=EventSeatsTable,AWS_ENDPOINT_URL=http://host.docker.internal:8000}" > /dev/null
        ;;
      "saga-process-payment")
        aws --endpoint-url $LOCALSTACK lambda create-function \
          --function-name "$name" --runtime provided.al2023 --architectures arm64 \
          --handler bootstrap --role arn:aws:iam::000000000000:role/dummy \
          --zip-file "fileb://$LAMBDAS_DIR/checkout-aws-lambda/checkout.zip" > /dev/null
        ;;
      "saga-create-ticket")
        aws --endpoint-url $LOCALSTACK lambda create-function \
          --function-name "$name" --runtime provided.al2023 --architectures arm64 \
          --handler bootstrap --role arn:aws:iam::000000000000:role/dummy \
          --zip-file "fileb://$LAMBDAS_DIR/tickets-aws-lambda/tickets.zip" \
          --environment "Variables={TICKETS_TABLE_NAME=TicketsTable,AWS_ENDPOINT_URL=http://host.docker.internal:8000}" > /dev/null
        ;;
    esac
    echo "  Created: $name"
  else
    echo "  Updating Lambda code: $name"
    case "$name" in
      "saga-reserve-seat"|"saga-release-seat")
        aws --endpoint-url $LOCALSTACK lambda update-function-code \
          --function-name "$name" \
          --zip-file "fileb://$LAMBDAS_DIR/seats-aws-lambda/seats.zip" > /dev/null
        ;;
      "saga-process-payment")
        aws --endpoint-url $LOCALSTACK lambda update-function-code \
          --function-name "$name" \
          --zip-file "fileb://$LAMBDAS_DIR/checkout-aws-lambda/checkout.zip" > /dev/null
        ;;
      "saga-create-ticket")
        aws --endpoint-url $LOCALSTACK lambda update-function-code \
          --function-name "$name" \
          --zip-file "fileb://$LAMBDAS_DIR/tickets-aws-lambda/tickets.zip" > /dev/null
        ;;
    esac
  fi
done

echo ""
echo ">>> [3/3] Creating/Updating Step Functions state machine..."

SAGA_STATUS=$(aws --endpoint-url $LOCALSTACK stepfunctions describe-state-machine \
  --state-machine-arn arn:aws:states:us-east-1:000000000000:stateMachine:checkout-saga 2>&1)

if echo "$SAGA_STATUS" | grep -q "StateMachineDoesNotExist"; then
  echo "  Creating state machine: checkout-saga"
  aws --endpoint-url $LOCALSTACK stepfunctions create-state-machine \
    --name checkout-saga \
    --role-arn arn:aws:iam::000000000000:role/dummy \
    --definition "file://$SAGA_DIR/checkout-saga.asl.json"
else
  echo "  Updating state machine: checkout-saga"
  aws --endpoint-url $LOCALSTACK stepfunctions update-state-machine \
    --state-machine-arn arn:aws:states:us-east-1:000000000000:stateMachine:checkout-saga \
    --definition "file://$SAGA_DIR/checkout-saga.asl.json"
fi

echo ""
echo "=== SAGA Setup Complete! ==="
echo "State Machine ARN: arn:aws:states:us-east-1:000000000000:stateMachine:checkout-saga"
echo "LocalStack Step Functions Endpoint: $LOCALSTACK"
