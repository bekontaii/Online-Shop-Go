# Online Shop API (Go backend)

A REST API for an online shop built in Go. The project follows Clean Architecture principles and utilizes Go's standard library for routing.

## Tech Stack

*   **Programming Language:** Go (version 1.25+)
*   **Database:** PostgreSQL (driver `github.com/lib/pq`)
*   **Authentication:** JWT (using `github.com/golang-jwt/jwt/v5`)
*   **Security:** Password hashing via Bcrypt (`golang.org/x/crypto/bcrypt`)

---

## Project Architecture

The project is split into logical layers and domains within the `internal` directory:

```text
├── cmd
│   └── main.go                 # Entry point, starts the server and initializes dependencies
├── internal
│   ├── auth                    # DTOs for registration and login
│   ├── cart                    # Customer cart management
│   ├── category                # Product category management (Admin only)
│   ├── middleware              # JWT middleware and context management
│   ├── order                   # Order creation and management (transactional checkout)
│   ├── product                 # Product catalog management (Seller/Admin)
│   ├── repository              # Shared repository interfaces (where applicable)
│   └── user                    # User registration, authorization, and profiles
├── pkg
│   ├── database                # PostgreSQL connection setup
│   └── jwt                     # JWT token generation and validation utilities
└── schema.sql                  # PostgreSQL database schema
```

Each business domain in `internal` contains:
*   `model.go` — entity data structures.
*   `dto.go` — API request and response data transfer objects.
*   `repository.go` — database abstraction interface.
*   `postgres_repository.go` — PostgreSQL repository implementation.
*   `service.go` — business logic and validation rules.
*   `handler.go` — HTTP handlers and routing (utilizing Go 1.22+ ServeMux).

---

## Installation & Setup

### 1. Database Setup
Ensure you have PostgreSQL installed and a database created (e.g., `onlineshop`).
Initialize the database tables using the `schema.sql` file:
```bash
psql -U postgres -d onlineshop -f schema.sql
```

### 2. Environment Variables
To enable JWT authentication, you must configure a secret key. In your terminal, run:

**Windows (PowerShell):**
```powershell
$env:JWT_SECRET="your_secret_key"
```

**Linux / macOS:**
```bash
export JWT_SECRET="your_secret_key"
```

### 3. Running the Server
Run the application using:
```bash
go run cmd/main.go
```
The server will start on port `8080` (`http://localhost:8080`).

---

## API Endpoints

### Authentication
*   `POST /api/register` — Register a new user (Roles: `user`, `seller`, `admin`).
*   `POST /api/login` — Login and receive a JWT token.

### Product Categories
*   `GET /api/categories` — Retrieve all categories.
*   `GET /api/categories/{id}` — Get a specific category.
*   `POST /api/categories` 🔒 — Create a category (Admin only).
*   `PATCH /api/categories/{id}` 🔒 — Update a category (Admin only).
*   `DELETE /api/categories/{id}` 🔒 — Delete a category (Admin only).

### Products (Catalog)
*   `GET /api/products` — Retrieve products with optional filtering by `category_id` and text search (`search`).
*   `GET /api/products/{id}` — Get details of a specific product.
*   `POST /api/products` 🔒 — Add a new product (Seller or Admin).
*   `PATCH /api/products/{id}` 🔒 — Edit a product (Product owner or Admin).
*   `DELETE /api/products/{id}` 🔒 — Remove a product (Product owner or Admin).

### Shopping Cart
*   `GET /api/cart` 🔒 — View the authorized customer's cart.
*   `POST /api/cart/add` 🔒 — Add an item to the cart.
*   `PATCH /api/cart/update` 🔒 — Update the quantity of an item in the cart.
*   `DELETE /api/cart/delete?product_id={id}` 🔒 — Remove an item from the cart.

### Orders & Checkout
*   `POST /api/orders/checkout` 🔒 — Checkout current cart items transactionally (verifies and decrements stock).
*   `GET /api/orders` 🔒 — Order history (customers see their own orders; admins and sellers view all orders).
*   `GET /api/orders/{id}` 🔒 — View details of a specific order.
*   `PATCH /api/orders/{id}/status` 🔒 — Update order status (Seller or Admin).

🔒 *Requires an `Authorization: Bearer <JWT_TOKEN>` header.*

---

## Testing Examples (cURL)

### Registration & Login
```bash
# 1. Register a new seller
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username": "john_seller", "email": "john@shop.com", "password": "supersecurepassword", "role": "seller"}'

# 2. Login to get token
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "john@shop.com", "password": "supersecurepassword"}'
```
*Save the returned JWT token to authorize subsequent API requests.*

### Adding a Product (Requires Seller or Admin role)
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Wireless Mouse", "description": "2.4 GHz optical mouse", "price": 25.50, "stock": 100, "category_id": 1}'
```

### Checkout Flow
```bash
# 1. Add product to cart
curl -X POST http://localhost:8080/api/cart/add \
  -H "Authorization: Bearer <USER_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "quantity": 3}'

# 2. Perform checkout
curl -X POST http://localhost:8080/api/orders/checkout \
  -H "Authorization: Bearer <USER_JWT_TOKEN>"
```
