# README.md

# Coupon System MVP

## Overview
This project implements a Coupon System as part of a medicine ordering platform. It provides backend logic for creating, managing, and validating coupons with a focus on correctness, modularity, and production readiness.

## Core Features
- **Admin Coupon Creation**: Create and manage coupon codes with attributes such as expiry date, usage type, applicable medicine IDs/categories, minimum order value, discount details, and usage limits.
- **Coupon Validation**: Validate coupons against constraints and ensure correct usage behavior.
- **Applicable Coupons**: Retrieve coupons applicable to a given cart and order.
- **Authentication Middleware**: Basic Bearer token authentication for protected endpoints.

## API Endpoints
- **POST /coupons**: Create a new coupon.
- **GET /coupons/applicable**: Retrieve applicable coupons based on cart items, order total, and timestamp (via query parameters).
- **POST /coupons/validate**: Validate a coupon code against the provided cart items, order total, and timestamp.

## Technical Features
- **Concurrency-aware design**: In-memory repository uses mutexes for thread safety.
- **Persistent storage**: Uses SQLite for storing coupons (can be extended to Postgres).
- **Request-scoped context**: All repository and service methods accept `context.Context`.
- **Caching**: Repository is wrapped with a TTL cache for coupon lookups.
- **OpenAPI documentation**: API is documented in `docs/swagger.yaml` and can be served via Swagger UI.
- **Dockerized**: Includes a Dockerfile for containerized deployment.

## Project Structure
```
coupon-system
├── cmd
│   └── server
│       └── main.go
├── internal
│   ├── api
│   │   ├── handlers
│   │   │   ├── coupon.go
│   │   │   └── routes.go
│   │   └── middleware
│   │       └── auth.go
│   ├── config
│   │   └── config.go
│   ├── models
│   │   └── coupon.go
│   ├── repository
│   │   └── coupon.go
│   ├── service
│   │   └── coupon.go
│   └── cache
│       └── cache.go
├── pkg
│   ├── validator
│   │   └── validator.go
│   └── errors
│       └── errors.go
├── docs
│   └── swagger.yaml
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

## Setup Instructions

1. **Clone the repository**
2. **Navigate to the project directory**
3. **Install dependencies**
   ```sh
   go mod tidy
   ```
4. **(Optional) Install SQLite CLI**  
   For database inspection, install SQLite CLI or use Docker as described below.
5. **Start the server**
   ```sh
   go run cmd/server/main.go
   ```
   The server will automatically create `coupons.db` in `cmd/server/` and set up the required table.

6. **(Optional) Inspect the database**
   If you have Docker:
   ```sh
   sudo docker run -it --rm -v "$PWD/cmd/server":/db nouchka/sqlite3 sqlite3 /db/coupons.db
   ```
   Then use `.tables` and `SELECT * FROM coupons;` to inspect data.

## Usage

- **Create a coupon**
  ```sh
  curl -X POST "http://localhost:8080/coupons" \
    -H "Content-Type: application/json" \
    -d '{
      "coupon_code": "SAVE20",
      "expiry_date": "2025-12-31T23:59:59Z",
      "usage_type": "one_time",
      "applicable_medicine_ids": ["med_123"],
      "applicable_categories": ["painkiller"],
      "min_order_value": 500,
      "discount_type": "fixed",
      "discount_value": 20,
      "max_usage_per_user": 1
    }'
  ```

- **Get applicable coupons**
  ```sh
  curl -G "http://localhost:8080/coupons/applicable" \
    --data-urlencode "order_total=700" \
    --data-urlencode "timestamp=2025-05-05T15:00:00Z" \
    --data-urlencode "cart_item_id=med_123" \
    --data-urlencode "cart_item_category=painkiller"
  ```

- **Validate a coupon**
  ```sh
  curl -X POST "http://localhost:8080/coupons/validate" \
    -H "Content-Type: application/json" \
    -d '{
      "coupon_code": "SAVE20",
      "cart_items": [
        { "id": "med_123", "category": "painkiller", "price": 100, "quantity": 2 }
      ],
      "order_total": 200,
      "timestamp": "2025-05-05T15:00:00Z"
    }'
  ```

## Swagger Documentation

- OpenAPI documentation is available at `docs/swagger.yaml`.
- To serve Swagger UI, integrate [swaggo/http-swagger](https://github.com/swaggo/http-swagger) in your router, or use an online Swagger editor.

## Docker

- **Build the Docker image**
  ```sh
  docker build -t coupon-system .
  ```
- **Run the container**
  ```sh
  docker run -p 8080:8080 coupon-system
  ```

## Notes

- The application is designed to handle concurrent requests safely.
- Caching is implemented for performance.
- The codebase is modular and ready for extension (e.g., Postgres support, advanced auth, etc.).

## Optional

- Deployed API URL can be provided upon request.# coupon-system
