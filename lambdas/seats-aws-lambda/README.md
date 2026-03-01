# Seats AWS Lambda

Component that handles seat availability checking, reservation constraints, and optimistic locking during high-concurrency scenarios before payment.

## Component Description

Manages seat status (available, reserved, sold) and guarantees atomic transactions in the DynamoDB table. Acts as a core executor in SAGA choreography.

## Technologies Used
- AWS Lambda
- Go
- DynamoDB Optimistic Locking
