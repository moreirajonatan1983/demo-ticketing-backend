#!/bin/bash
# start_all.sh - Inicia todos los servicios de la Ticketera localmente (Full Stack)

# Colores para el log
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

BASE_DIR="/Users/jonatandanielmoreira/developer/proyectos/demo"
SCRIPTS_DIR="$BASE_DIR/demo-ticketing-backend/scripts"

echo -e "${BLUE}>>> Iniciando entorno local de Ticketera Cloud (API & Web)...${NC}"

# Verificación de Docker (SAM lo necesita)
DOCKER_RUNNING=true
if ! docker info >/dev/null 2>&1; then
    echo -e "${RED}AVISO: Docker no está corriendo. Las APIs Backend fallarán, pero la Web App React levantará modo Mock.${NC}"
    DOCKER_RUNNING=false
fi

# 1. Construcción de Binarios
echo -e "${GREEN}>>> [1/3] Construyendo Backend (Go/Node) con AWS SAM...${NC}"
cd "$BASE_DIR/demo-ticketing-backend/lambdas/events-aws-lambda" 
if [ "$DOCKER_RUNNING" = true ] && [ -f "template.yaml" ]; then
    sam build > build-backend-events.log 2>&1
else
    echo -e "${RED}Aviso: template.yaml no encontrado o Docker apagado en events. Omitiendo SAM build...${NC}"
fi

cd "$BASE_DIR/demo-ticketing-backend/lambdas/checkout-aws-lambda" 
if [ "$DOCKER_RUNNING" = true ] && [ -f "template.yaml" ]; then
    sam build > build-backend-checkout.log 2>&1
else
    echo -e "${RED}Aviso: template.yaml no encontrado o Docker apagado en checkout. Omitiendo SAM build...${NC}"
fi

cd "$BASE_DIR/demo-ticketing-backend/lambdas/seats-aws-lambda" 
if [ "$DOCKER_RUNNING" = true ] && [ -f "template.yaml" ]; then
    sam build > build-backend-seats.log 2>&1
else
    echo -e "${RED}Aviso: template.yaml no encontrado o Docker apagado en seats. Omitiendo SAM build...${NC}"
fi

cd "$BASE_DIR/demo-ticketing-backend/lambdas/tickets-aws-lambda" 
if [ "$DOCKER_RUNNING" = true ] && [ -f "template.yaml" ]; then
    sam build > build-backend-tickets.log 2>&1
else
    echo -e "${RED}Aviso: template.yaml no encontrado o Docker apagado en tickets. Omitiendo SAM build...${NC}"
fi

cd "$BASE_DIR/demo-ticketing-backend/lambdas/shows-aws-lambda" 
if [ "$DOCKER_RUNNING" = true ] && [ -f "template.yaml" ]; then
    sam build > build-backend-shows.log 2>&1
else
    echo -e "${RED}Aviso: template.yaml no encontrado o Docker apagado en shows. Omitiendo SAM build...${NC}"
fi

# 2. Levantar APIs en Segundo Plano
echo -e "${BLUE}>>> [2/3] Levantando simuladores de API Gateway...${NC}"

# Backend Transaccional - Events (Puerto 3000)
cd "$BASE_DIR/demo-ticketing-backend/lambdas/events-aws-lambda"
if [ "$DOCKER_RUNNING" = true ] && [ -f "template.yaml" ]; then
    AWS_ACCESS_KEY_ID="test" AWS_SECRET_ACCESS_KEY="test" AWS_REGION="us-east-1" nohup sam local start-api --host 0.0.0.0 --env-vars "$BASE_DIR/demo-ticketing-backend/env.json" --port 3000 > sam-backend-events.log 2>&1 &
    echo $! > "$SCRIPTS_DIR/.backend_api.pid"
else
    echo -e "${RED}Aviso: Omitiendo start-api de events (Puerto 3000).${NC}"
fi

# Backend Transaccional - Checkout (Puerto 3004)
cd "$BASE_DIR/demo-ticketing-backend/lambdas/checkout-aws-lambda"
if [ "$DOCKER_RUNNING" = true ] && [ -f "template.yaml" ]; then
    AWS_ACCESS_KEY_ID="test" AWS_SECRET_ACCESS_KEY="test" AWS_REGION="us-east-1" nohup sam local start-api --host 0.0.0.0 --env-vars "$BASE_DIR/demo-ticketing-backend/env.json" --port 3004 > sam-backend-checkout.log 2>&1 &
    echo $! > "$SCRIPTS_DIR/.checkout_api.pid"
else
    echo -e "${RED}Aviso: Omitiendo start-api de checkout (Puerto 3004).${NC}"
fi

# Backend Transaccional - Seats (Puerto 3005)
cd "$BASE_DIR/demo-ticketing-backend/lambdas/seats-aws-lambda"
if [ "$DOCKER_RUNNING" = true ] && [ -f "template.yaml" ]; then
    AWS_ACCESS_KEY_ID="test" AWS_SECRET_ACCESS_KEY="test" AWS_REGION="us-east-1" nohup sam local start-api --host 0.0.0.0 --env-vars "$BASE_DIR/demo-ticketing-backend/env.json" --port 3005 > sam-backend-seats.log 2>&1 &
    echo $! > "$SCRIPTS_DIR/.seats_api.pid"
else
    echo -e "${RED}Aviso: Omitiendo start-api de seats (Puerto 3005).${NC}"
fi

# Backend Transaccional - Tickets (Puerto 3006)
cd "$BASE_DIR/demo-ticketing-backend/lambdas/tickets-aws-lambda"
if [ "$DOCKER_RUNNING" = true ] && [ -f "template.yaml" ]; then
    AWS_ACCESS_KEY_ID="test" AWS_SECRET_ACCESS_KEY="test" AWS_REGION="us-east-1" nohup sam local start-api --host 0.0.0.0 --env-vars "$BASE_DIR/demo-ticketing-backend/env.json" --port 3006 > sam-backend-tickets.log 2>&1 &
    echo $! > "$SCRIPTS_DIR/.tickets_api.pid"
else
    echo -e "${RED}Aviso: Omitiendo start-api de tickets (Puerto 3006).${NC}"
fi

# Backend Transaccional - Shows (Puerto 3007)
cd "$BASE_DIR/demo-ticketing-backend/lambdas/shows-aws-lambda"
if [ "$DOCKER_RUNNING" = true ] && [ -f "template.yaml" ]; then
    AWS_ACCESS_KEY_ID="test" AWS_SECRET_ACCESS_KEY="test" AWS_REGION="us-east-1" nohup sam local start-api --host 0.0.0.0 --env-vars "$BASE_DIR/demo-ticketing-backend/env.json" --port 3007 > sam-backend-shows.log 2>&1 &
    echo $! > "$SCRIPTS_DIR/.shows_api.pid"
else
    echo -e "${RED}Aviso: Omitiendo start-api de shows (Puerto 3007).${NC}"
fi

# Auth Backend (Puerto 3003)
cd "$BASE_DIR/demo-ticketing-backend/lambdas/auth-aws-lambda"
if [ "$DOCKER_RUNNING" = true ] && [ -f "template.yaml" ]; then
    sam build > build-backend-auth.log 2>&1
    AWS_ACCESS_KEY_ID="test" AWS_SECRET_ACCESS_KEY="test" AWS_REGION="us-east-1" nohup sam local start-api --host 0.0.0.0 --env-vars "$BASE_DIR/demo-ticketing-backend/env.json" --port 3003 > sam-auth.log 2>&1 &
    echo $! > "$SCRIPTS_DIR/.auth_api.pid"
    echo -e "${GREEN}  ✅ Auth API levantando en puerto 3003${NC}"
else
    echo -e "${RED}Aviso: Omitiendo start-api de auth (Puerto 3003).${NC}"
fi

# Worker/Java Backend (Puerto 3002 - Reservado para emular Playground)
cd "$BASE_DIR/demo-ticketing-worker" 2>/dev/null || cd "$BASE_DIR"
if [ "$DOCKER_RUNNING" = true ] && [ -f "template.yaml" ]; then
    nohup sam local start-api --host 0.0.0.0 --docker-network demo-ticketing-backend_default --port 3002 > sam-worker.log 2>&1 &
    echo $! > "$SCRIPTS_DIR/.worker_api.pid"
else
    echo -e "${RED}Aviso: Omitiendo start-api de worker (Puerto 3002).${NC}"
fi

# 3. Iniciar Frontend (Web)
echo -e "${GREEN}>>> [3/3] Iniciando React Web App (Vite)...${NC}"
cd "$BASE_DIR/demo-ticketing-web"
nohup npm run dev -- --port 3001 > web.log 2>&1 &
echo $! > "$SCRIPTS_DIR/.web.pid"

echo -e "${BLUE}================================================================${NC}"
echo -e "${GREEN}¡Ticketera iniciada en background!${NC}"
echo -e "Endpoints Locales:"
echo -e " 🚀 Frontend (Vite):       http://localhost:3001"
echo -e " 📦 Backend Core (Events): http://localhost:3000"
echo -e " 📦 Backend Core (CheckOut): http://localhost:3004"
echo -e " 📦 Backend Core (Seats):  http://localhost:3005"
echo -e " 📦 Backend Core (Tickets):http://localhost:3006"
echo -e " 📦 Backend Core (Shows):  http://localhost:3007"
echo -e " 🔐 Auth Backend (Mock):   http://localhost:3003"
echo -e " ☕ Worker API (Mock):     http://localhost:3002"
echo -e "${BLUE}================================================================${NC}"
echo -e "Logs disponibles en directorios de servicios."
echo -e "Usa './scripts/stop_all.sh' para cerrar todos los procesos de forma limpia."
