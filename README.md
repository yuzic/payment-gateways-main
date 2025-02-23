# Payment Gateway Integration

### OpenApi Specification

```
./openapi.yaml

```

## Task Overview

Folder Structure
```
payment-gateway-service/
├── openapi/           # OpenApi entry points
├── cmd/               # Application entry points
├── db/                # Database operations
├── internal/          # Internal services and models
│   ├── api/           # API handlers
│   ├── kafka/         # Kafka producers
│   ├── models/        # Request/response and database models
│   ├── services/      # Core business logic
│   ├── util/          # Utility functions
├── docs/              # API documentation (OpenAPI)
```


### Endpoints

```Deposit Endpoint
URL: /deposit
Method: POST
Description: Processes deposit transactions.
Request Body Example:

{
    "amount": 100.00,
    "user_id": 1,
    "currency": "EUR"
}


```

```Deposit Endpoint
URL: /withdrawal
Method: POST
Description: Processes withdrawal transactions.
Request Body Example:

{
    "amount": 100.00,
    "user_id": 1,
    "currency": "EUR"
}

```

```
Callback Endpoint

URL: /callback
Method: GET
Description: Handles asynchronous callbacks from payment gateways to update transaction statuses.
Query Parameters:
id: Transaction ID
status: New transaction status (e.g., done, failed, pending)
gateway: Gateway identifier
Response: Returns a confirmation message after updating the transaction.
```




### How to Get Started

1. **Clone the Repository:**


2. **Setup Docker:**
    Docker is configured to run PostgreSQL, Kafka, and Redis. Use the following command to start all the services:

    ```bash
    docker-compose up -d
    ```

    This will start:
    - PostgreSQL on port `5432`
    - Kafka on ports `9092` and `9093`
    - Redis on port `6379`
    - Application on port `8080`

3. **Database Migration:**
    The migration file `db/init.sql` is already provided. Once the Docker services are up and running, the database will be initialized automatically, and the tables will be created.


