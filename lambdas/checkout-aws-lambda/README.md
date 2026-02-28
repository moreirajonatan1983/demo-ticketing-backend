# Checkout AWS Lambda

Este microservicio es el encargado de orquestar el flujo de compras. Recibe la orden de pedido, procesa el pago (mockeado) e interactúa bloqueando o reservando temporalmente los asientos adquiridos hasta proveer al usuario la confirmación.

## Arquitectura
Implementado en **Go** utilizando **Arquitectura Hexagonal**:
- **Domain**: Define el modelo del payload transaccional `CheckoutRequest` y la salida `CheckoutResponse`.
- **Ports**: Interfaces de entrada para los casos de uso (`CheckoutService`).
- **Services**: Lógica de negocio core (en proceso de desarrollo interactivo - async).
- **Adapters**: 
  - `handlers`: Adaptador HTTP / API Gateway que implementa Swagger (OpenAPI).

## Endpoints Principales
- `POST /checkout`: Recibe el ID de evento, el medio de pago y las butacas y devuelve una orden asincrónica.
- `GET /swagger.json`: Especificación OpenAPI autogenerada.

## Ejecución Local
Este proyecto está configurado para ejecutarse con **AWS SAM CLI**.

```bash
cd checkout-aws-lambda
sam local start-api -p 3004 --env-vars ../../env.json
```
