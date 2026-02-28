# Tickets AWS Lambda

Este microservicio se encarga de la visualización, entrega y control de acceso de los tickets emitidos a nombre de cada usuario de forma segura. Provee la data necesaria para emitir un e-ticket digital (QR Code).

## Arquitectura
Implementado en **Go** utilizando **Arquitectura Hexagonal**:
- **Domain**: Define el modelo `Ticket`.
- **Ports**: Interfaces para los casos de uso (`TicketService`) y la base de datos (`TicketRepository`).
- **Services**: Lógica de negocio core.
- **Adapters**: 
  - `handlers`: Adaptador HTTP / API Gateway que implementa Swagger (OpenAPI).
  - `repositories`: Adaptador **DynamoDB** para lectura de tickets emitidos (TicketsTable).

## Endpoints Principales
- `GET /tickets/me`: Retorna los detalles de las entradas adquiridas correspondientes al usuario logueado en base a su Token.
- `GET /swagger.json`: Especificación OpenAPI autogenerada.

## Ejecución Local
Este proyecto está configurado para ejecutarse con **AWS SAM CLI** y **DynamoDB Local**.
Asegurate de que DynamoDB Local esté corriendo y las tablas estén creadas.

```bash
cd tickets-aws-lambda
sam local start-api -p 3006 --env-vars ../../env.json
```
