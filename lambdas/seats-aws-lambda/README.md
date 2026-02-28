# Seats AWS Lambda

Este microservicio se encarga de gestionar la disponibilidad y el mapa del estadio (Asientos) para un evento en particular.

## Arquitectura
Implementado en **Go** utilizando **Arquitectura Hexagonal**:
- **Domain**: Define el modelo `Seat`.
- **Ports**: Interfaces para los casos de uso (`SeatService`) y base de datos (`SeatRepository`).
- **Services**: Lógica de negocio core.
- **Adapters**: 
  - `handlers`: Adaptador HTTP / API Gateway que implementa Swagger (OpenAPI).
  - `repositories`: Adaptador **DynamoDB** para persistencia de datos (EventSeatsTable).

## Endpoints Principales
- `GET /events/{eventId}/seats`: Retorna la lista de asientos correspondientes a un evento indicando su estado `available`, `occupied` o `processing`.
- `GET /swagger.json`: Especificación OpenAPI autogenerada.

## Ejecución Local
Este proyecto está configurado para ejecutarse con **AWS SAM CLI** y **DynamoDB Local**.
Asegurate de que DynamoDB Local esté corriendo y las tablas estén creadas.

```bash
cd seats-aws-lambda
sam local start-api -p 3005 --env-vars ../../env.json
```
