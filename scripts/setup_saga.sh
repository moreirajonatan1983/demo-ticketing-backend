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
          --environment "Variables={SEATS_TABLE_NAME=EventSeatsTable,AWS_ENDPOINT_URL=http://host.docker.internal:4566}" > /dev/null
        ;;
      "saga-release-seat")
        aws --endpoint-url $LOCALSTACK lambda create-function \
          --function-name "$name" --runtime provided.al2023 --architectures arm64 \
          --handler bootstrap --role arn:aws:iam::000000000000:role/dummy \
          --zip-file "fileb://$LAMBDAS_DIR/seats-aws-lambda/seats.zip" \
          --environment "Variables={SEATS_TABLE_NAME=EventSeatsTable,AWS_ENDPOINT_URL=http://host.docker.internal:4566}" > /dev/null
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
          --environment "Variables={TICKETS_TABLE_NAME=TicketsTable,AWS_ENDPOINT_URL=http://host.docker.internal:4566}" > /dev/null
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
        aws --endpoint-url $LOCALSTACK lambda update-function-configuration \
          --function-name "$name" \
          --environment "Variables={SEATS_TABLE_NAME=EventSeatsTable,AWS_ENDPOINT_URL=http://host.docker.internal:4566}" > /dev/null
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
        aws --endpoint-url $LOCALSTACK lambda update-function-configuration \
          --function-name "$name" \
          --environment "Variables={TICKETS_TABLE_NAME=TicketsTable,AWS_ENDPOINT_URL=http://host.docker.internal:4566}" > /dev/null
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

# ─────────────────────────────────────────────────────────────
# [4/4] Crear API Key + Usage Plan en LocalStack (Opción B Auth)
# ─────────────────────────────────────────────────────────────
echo ""
echo ">>> [4/4] Configurando API Key + Usage Plan en LocalStack..."

# Nombre fijo para idempotencia
API_KEY_NAME="ticketera-web-client"
USAGE_PLAN_NAME="ticketera-default-plan"

# Verificar si ya existe la API Key
EXISTING_KEY=$(aws --endpoint-url $LOCALSTACK apigateway get-api-keys \
  --name-query "$API_KEY_NAME" \
  --include-values 2>/dev/null | grep '"value"' | head -1 | sed 's/.*"value": "\(.*\)".*/\1/')

if [ -n "$EXISTING_KEY" ]; then
  echo "  API Key ya existe: $API_KEY_NAME"
  API_KEY_VALUE="$EXISTING_KEY"
else
  echo "  Creando API Key: $API_KEY_NAME"
  CREATE_RESULT=$(aws --endpoint-url $LOCALSTACK apigateway create-api-key \
    --name "$API_KEY_NAME" \
    --description "Ticketera Web Client - Authorized callers only" \
    --enabled \
    --generate-distinct-id 2>/dev/null)
  
  API_KEY_ID=$(echo "$CREATE_RESULT" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])" 2>/dev/null)
  API_KEY_VALUE=$(echo "$CREATE_RESULT" | python3 -c "import sys,json; print(json.load(sys.stdin)['value'])" 2>/dev/null)
  
  echo "  API Key creada: ID=$API_KEY_ID"
fi

# Crear Usage Plan si no existe
EXISTING_PLAN=$(aws --endpoint-url $LOCALSTACK apigateway get-usage-plans 2>/dev/null | \
  python3 -c "import sys,json; plans=json.load(sys.stdin).get('items',[]); [print(p['id']) for p in plans if p['name']=='$USAGE_PLAN_NAME']" 2>/dev/null)

if [ -z "$EXISTING_PLAN" ]; then
  echo "  Creando Usage Plan: $USAGE_PLAN_NAME"
  PLAN_RESULT=$(aws --endpoint-url $LOCALSTACK apigateway create-usage-plan \
    --name "$USAGE_PLAN_NAME" \
    --description "Plan de uso para clientes autorizados de Ticketera" \
    --throttle burstLimit=100,rateLimit=50 \
    --quota limit=10000,offset=0,period=MONTH 2>/dev/null)
  PLAN_ID=$(echo "$PLAN_RESULT" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])" 2>/dev/null)
  
  if [ -n "$API_KEY_ID" ] && [ -n "$PLAN_ID" ]; then
    aws --endpoint-url $LOCALSTACK apigateway create-usage-plan-key \
      --usage-plan-id "$PLAN_ID" \
      --key-id "$API_KEY_ID" \
      --key-type "API_KEY" 2>/dev/null
    echo "  API Key asociada al Usage Plan"
  fi
else
  echo "  Usage Plan ya existe: $USAGE_PLAN_NAME"
fi

# Escribir la API Key al .env.local del frontend (si existe el directorio)
WEB_DIR="$SCRIPT_DIR/../../demo-ticketing-web"
if [ -d "$WEB_DIR" ] && [ -n "$API_KEY_VALUE" ]; then
  ENV_FILE="$WEB_DIR/.env.local"
  # Actualizar o crear la variable VITE_API_KEY
  if [ -f "$ENV_FILE" ]; then
    # Reemplazar si existe
    if grep -q "VITE_API_KEY" "$ENV_FILE"; then
      sed -i.bak "s|^VITE_API_KEY=.*|VITE_API_KEY=$API_KEY_VALUE|" "$ENV_FILE" && rm -f "$ENV_FILE.bak"
    else
      echo "VITE_API_KEY=$API_KEY_VALUE" >> "$ENV_FILE"
    fi
  else
    echo "VITE_API_KEY=$API_KEY_VALUE" > "$ENV_FILE"
  fi
  echo "  API Key escrita en $ENV_FILE"
fi

echo ""
echo "  API Key Value: ${API_KEY_VALUE:-'(ya existente, verificar .env.local)'}"
echo ""
echo "  Para Postman: agregar header   x-api-key: $API_KEY_VALUE"
echo "  Para Postman: agregar header   Authorization: Bearer <id_token_from_login>"
echo ""

