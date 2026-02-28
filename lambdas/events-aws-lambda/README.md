# Events AWS Lambda

Este microservicio se encarga de gestionar el listado de eventos de la plataforma (Cartelera).

## Arquitectura
Implementado en **Go** utilizando **Arquitectura Hexagonal**:
- **Domain**: Define el modelo `Event`.
- **Ports**: Interfaces para los casos de uso (`EventService`) y base de datos (`EventRepository`).
- **Services**: Lógica de negocio core.
- **Adapters**: 
  - `handlers`: Adaptador HTTP / API Gateway que implementa Swagger (OpenAPI).
  - `repositories`: Adaptador **DynamoDB** para persistencia de datos (EventsTable).

## Endpoints Principales
- `GET /events`: Retorna todos los eventos disponibles.
- `GET /events/{id}`: Retorna el detalle de un evento específico.
- `GET /swagger.json`: Especificación OpenAPI autogenerada.

## Ejecución Local
Este proyecto está configurado para ejecutarse con **AWS SAM CLI** y **DynamoDB Local**.
Asegurate de que DynamoDB Local esté corriendo y las tablas estén creadas.

```bash
cd events-aws-lambda
sam local start-api -p 3000 --env-vars ../../env.json
```
