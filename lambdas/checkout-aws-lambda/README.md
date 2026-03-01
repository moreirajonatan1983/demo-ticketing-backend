# Checkout AWS Lambda

This is the Checkout component responsible for coordinating the SAGA checkout process and payment validation for the demo ticketing application.

## Component Description

Handles payment processing and initiating saga orchestrations for tickets when users finalize their orders.

## Technologies Used
- AWS Lambda
- Go
- CQRS / Saga orchestration (via Step Functions)
