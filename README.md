# 📚 Library System — Microservices

A microservices-based **library management system** built with **Go, Gin, and PostgreSQL**. Members can browse books, borrow and return them, and the system tracks fines — split into independent services behind an API gateway.

## 🏗️ Architecture

The system follows a **microservices architecture** with an **API Gateway**:

```
                    Client (Postman / frontend)
                              │
                    ┌─────────▼─────────┐
                    │    API Gateway    │  :8080   (reverse proxy)
                    └────┬─────────┬────┘
             /api/loan   │         │  /api/catalog
                 ┌───────▼──┐   ┌──▼────────┐
                 │   Loan   │◄──│  Catalog  │   (loan calls catalog via REST)
                 │ Service  │   │  Service  │
                 │  :8082   │   │  :8081    │
                 └─────┬────┘   └─────┬─────┘
                    loan_db        catalog_db     (database per service)
```

**Key patterns:**
- **Microservices** — independent, separately deployable services
- **API Gateway** — single entry point (reverse proxy); routes by URL prefix
- **Database-per-service** — each service owns its PostgreSQL database (no shared DB)
- **REST over HTTP** — services communicate synchronously
- **Stateless JWT authentication** — issued by the loan service, verified by both with a shared secret

## 🧩 Services

| Service | Port | Responsibility | Database |
|---|---|---|---|
| **Gateway** | 8080 | Single entry point; routes requests to services | — |
| **Loan Service** | 8082 | Members, authentication, borrowing, returns, fines | `loan_db` |
| **Catalog Service** | 8081 | Books (browse + librarian management) | `catalog_db` |

## 🛠️ Tech Stack
- **Go** + **Gin** — HTTP framework
- **PostgreSQL** + **GORM** — database & ORM (auto-migration)
- **JWT** (`golang-jwt`) — stateless authentication
- **Swagger** (`swaggo`) — interactive API documentation
- **godotenv** — configuration via environment variables

## 📁 Project Structure

```
library-system/
├── gateway/                # API gateway (reverse proxy)
├── loan-service/           # members, auth, loans, fines
│   ├── main.go
│   └── internal/
│       ├── controllers/    # request handlers
│       ├── models/         # Member, Loan, Fine
│       ├── middleware/     # JWT auth
│       ├── routes/         # URL → controller wiring
│       ├── database/       # DB connection + librarian seed
│       └── clients/        # HTTP client that calls catalog
├── catalog-service/        # books
└── README.md
```

## ✅ Prerequisites
- **Go** 1.21+
- **PostgreSQL** 15+ (with pgAdmin)

## ⚙️ Setup

**1. Clone the repository**
```bash
git clone https://github.com/Yostena16/library-system.git
cd library-system
```

**2. Create the databases** (in pgAdmin)
- `loan_db`
- `catalog_db`

**3. Create a `.env` file in each service** (values are read from the environment, never hardcoded)

`loan-service/.env`
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=loan_db
JWT_SECRET=your-shared-secret
JWT_EXPIRY_HOURS=24
CATALOG_SERVICE_URL=http://localhost:8081
LIBRARIAN_EMAIL=librarian@library.com
LIBRARIAN_PASSWORD=librarian123
```

`catalog-service/.env`
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=catalog_db
JWT_SECRET=your-shared-secret
```

`gateway/.env`
```env
CATALOG_SERVICE_URL=http://localhost:8081
LOAN_SERVICE_URL=http://localhost:8082
```

> ⚠️ **`JWT_SECRET` must be identical** in the loan and catalog services so both can verify the same tokens.

**4. Install dependencies** (run inside each service folder)
```bash
go mod tidy
```

## ▶️ Running

Start each service in its own terminal:

```bash
# Terminal 1 — Catalog
cd catalog-service && go run main.go   # http://localhost:8081

# Terminal 2 — Loan
cd loan-service && go run main.go       # http://localhost:8082

# Terminal 3 — Gateway
cd gateway && go run main.go            # http://localhost:8080
```

All client requests go through the **gateway** at `http://localhost:8080`.

## 📖 API Documentation (Swagger)

- **Loan Service:** http://localhost:8082/swagger/index.html
- **Catalog Service:** http://localhost:8081/swagger/index.html

## 🔐 Authentication & Roles

Authentication uses **JWT**. Log in to receive a token, then send it on protected routes as:
```
Authorization: Bearer <token>
```

| Role | Permissions |
|---|---|
| 🌐 **Public** | Browse books, register, login |
| 👤 **Member** | Borrow books, view own loans & fines |
| 🧑‍💼 **Librarian** | Manage books, process returns |

- Registration always creates a **member**.
- A **librarian** account is **seeded automatically on startup** from the loan service's `.env` (it is never created via public registration).

## 🌐 Main Endpoints (via gateway)

**Auth** (public)
| Method | Endpoint |
|---|---|
| POST | `/api/loan/auth/register` |
| POST | `/api/loan/auth/login` |

**Members / Loans / Fines** (require token)
| Method | Endpoint | Who |
|---|---|---|
| GET | `/api/loan/members/me` | member |
| POST | `/api/loan/loans` | member (borrow) |
| POST | `/api/loan/loans/:id/return` | librarian |
| GET | `/api/loan/loans` | member |
| GET | `/api/loan/fines` | member |

**Catalog**
| Method | Endpoint | Who |
|---|---|---|
| GET | `/api/catalog/books` | public |
| GET | `/api/catalog/books/:id/availability` | public |
| POST/PUT/DELETE | `/api/catalog/books` | librarian |

## 🔗 How the Services Communicate

When a member borrows a book, the **Loan Service** makes an HTTP call to the **Catalog Service** (`GET /books/:id/availability`) to check availability before creating the loan:

```
borrow request → Loan Service → (HTTP) → Catalog Service → "available?"
              → creates the loan in loan_db
```

This is the core microservice interaction — two independent services cooperating over REST, each with its own database.

---

Built as a learning project demonstrating **microservices, an API gateway, JWT authentication with roles, database-per-service, and inter-service communication** in Go.
