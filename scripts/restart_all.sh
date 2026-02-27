#!/bin/bash
# restart_all.sh - Limpia los procesos y reinicia todo el entorno Ticketera

BLUE='\033[0;34m'
NC='\033[0m'

SCRIPTS_DIR="/Users/jonatandanielmoreira/developer/proyectos/demo/demo-ticketing-backend/scripts"

echo -e "${BLUE}>>> Reiniciando el entorno general de la Ticketera...${NC}"

bash "$SCRIPTS_DIR/stop_all.sh"
sleep 2
bash "$SCRIPTS_DIR/start_all.sh"
