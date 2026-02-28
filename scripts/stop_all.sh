#!/bin/bash
# stop_all.sh - Detiene de forma segura todos los servicios del proyecto Ticketera

RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

BASE_DIR="/Users/jonatandanielmoreira/developer/proyectos/demo"
SCRIPTS_DIR="$BASE_DIR/demo-ticketing-backend/scripts"

echo -e "${BLUE}>>> Deteniendo el entorno local de Ticketera Cloud...${NC}"

# Matar procesos registrados por PID
for service in backend_api checkout_api auth_api worker_api web; do
    pid_file="$SCRIPTS_DIR/.$service.pid"
    if [ -f "$pid_file" ]; then
        pid=$(cat "$pid_file")
        if ps -p $pid > /dev/null; then
            echo -e "Cerrando $service (PID: $pid)..."
            kill -9 $pid 2>/dev/null
        fi
        rm "$pid_file"
    fi
done

# Limpieza global de procesos huérfanos de SAM y Vite
echo -e "${RED}Limpiando procesos remanentes y contenedores muertos...${NC}"
pkill -f "sam local start-api" 2>/dev/null
pkill -f "vite" 2>/dev/null

echo -e "${BLUE}>>> Todos los servicios han sido detenidos exitosamente.${NC}"
