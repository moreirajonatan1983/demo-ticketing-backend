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
cd "$BASE_DIR/demo-ticketing-backend" 
if [ "$DOCKER_RUNNING" = true ] && [ -f "template.yaml" ]; then
    sam build > build-backend.log 2>&1
else
    echo -e "${RED}Aviso: template.yaml no encontrado o Docker apagado. Omitiendo SAM build...${NC}"
fi

# 2. Levantar APIs en Segundo Plano
echo -e "${BLUE}>>> [2/3] Levantando simuladores de API Gateway...${NC}"

# Backend Transaccional (Puerto 3000)
cd "$BASE_DIR/demo-ticketing-backend"
if [ "$DOCKER_RUNNING" = true ] && [ -f "template.yaml" ]; then
    nohup sam local start-api --port 3000 > sam-backend.log 2>&1 &
    echo $! > "$SCRIPTS_DIR/.backend_api.pid"
else
    echo -e "${RED}Aviso: Omitiendo start-api de backend.${NC}"
fi

# 3. Iniciar Frontend (Web)
echo -e "${GREEN}>>> [3/3] Iniciando React Web App (Vite)...${NC}"
cd "$BASE_DIR/demo-ticketing-web"
nohup npm run dev -- --port 5173 > web.log 2>&1 &
echo $! > "$SCRIPTS_DIR/.web.pid"

echo -e "${BLUE}================================================================${NC}"
echo -e "${GREEN}¡Ticketera iniciada en background!${NC}"
echo -e "Endpoints Locales:"
echo -e " 🚀 Frontend (Vite):       http://localhost:5173"
echo -e " 📦 Backend Core (SAM):    http://localhost:3000"
echo -e "${BLUE}================================================================${NC}"
echo -e "Logs disponibles: demo-ticketing-backend/sam-backend.log y demo-ticketing-web/web.log"
echo -e "Usa './scripts/stop_all.sh' para cerrar todos los procesos de forma limpia."
