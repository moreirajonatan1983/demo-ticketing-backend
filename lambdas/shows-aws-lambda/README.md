# Shows AWS Lambda

Este microservicio se encarga de gestionar las funciones (fechas) disponibles de un Evento.

## Arquitectura
Implementado en **Go** utilizando **Arquitectura Hexagonal**.

## Endpoints Principales
- `GET /events/{eventId}/shows`: Retorna la lista de fechas correspondientes a un evento indicando su estado `available` o `soldout`.
- `GET /swagger.json`: Especificación OpenAPI autogenerada.

## Ejecución Local
Este proyecto está configurado para ejecutarse con **AWS SAM CLI** y **DynamoDB Local**.
Asegurate de que DynamoDB Local esté corriendo y las tablas estén creadas.

```bash
cd shows-aws-lambda
sam local start-api -p 3007 --env-vars ../../env.json
```
