# Ticketera Cloud - Core System & Workers

Núcleo transaccional de la Ticketera, el cual soporta los procesos críticos (Flash Crowds, compra y reserva, envío de tickets y reportes).

## Funciones Principales y Flujos

1.  **Orquestación de la Compra de la Entrada (Saga Pattern)**
    *   Gestionado mediante **AWS Step Functions**. Coordina la reserva atómica en DynamoDB, procesa el pago de forma condicional (usando el patrón **Circuit Breaker** frente a una pasarela simulada) y confirma on emite error revirtiendo el inventario.
2.  **Manejo Asíncrono de Notificaciones (Pub/Sub Pattern)**
    *   Uso intensivo de **Amazon EventBridge**. Una vez emitida la orden (Entrada Vendida), esto detona el envío de correos (Amazon SES) mediante Lambdas y métricas incrementales en CloudWatch desacopladas de la transacción principal del usuario.
3.  **Procesos Batch & Reportes Pesados (Kubernetes - EKS)**
    *   Contenedores haciendo "polling" desde **Amazon SQS** asumen cargas batch, tales como la generación de PDFs masivos con los tickets, la ingesta y manipulación enorme de datos para reportes analíticos para productoras de eventos, delegados lejos de las Lambdas para evadir timeouts y límites de CPU de la capa Serverless.
4.  **Trazabilidad Top-Tier**
    *   El Core inyecta correlación de trazas usando **AWS X-Ray** globalmente (midiendo cuánto tarda la base Dynamo frente a la subida al S3) y métricas alarmadas en tiempo real **(CloudWatch)**.
