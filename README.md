# Library System

A microservices-based library management system built with Go, Gin, and PostgreSQL.

## Architecture

- **gateway** (:8080) — single entry point that routes requests to the services
- **catalog-service** (:8081) — books, authors, copies; owns `catalog_db`
- **loan-service** (:8082) — members, auth, borrowing & fines; owns `loan_db`

Each service is independent, has its own database, and talks to others over HTTP.

## Tech stack

Go · Gin · PostgreSQL · JWT · Swagger · Docker Compose

## Status

🚧 Under construction.