# Chirpy-project

This is was guided course project from Botdev: "HTTP servers with Go".

## Overview

Chirpy is a simple Go web server and REST API that demonstrates:
- building an HTTP server from scratch with `net/http`
- routing requests with `http.ServeMux`
- user management, authentication, and JWT session handling
- PostgreSQL database access using generated SQL code (`sqlc`)
- basic metrics, health-checks, and admin controls
- static file serving for a front-end under `/app/`

## Project Structure

- `main.go` - application entry point
- `server/` - HTTP handler and server setup
- `internal/dbman/` - database connection and repository functions
- `internal/auth/` - password hashing, JWT creation/validation, and token helpers
- `internal/database/` - generated SQL models and queries
- `sql/` - schema and query definitions for `sqlc`
- `app/` - static front-end assets and entrypoint

## Features

- `POST /api/users` - create a new user
- `POST /api/login` - authenticate and receive JWT + refresh token
- `PUT /api/users` - update user credentials
- `POST /api/chirps` - post a new chirp (message)
- `GET /api/chirps` - list all chirps
- `GET /api/chirps/{id}` - retrieve a chirp by ID
- `DELETE /api/chirps/{id}` - delete a chirp by ID
- `POST /api/refresh` - refresh access tokens
- `POST /api/revoke` - revoke a refresh token
- `POST /api/polka/webhooks` - custom webhook endpoint
- `GET /admin/healthz` - server health check
- `GET /admin/metrics` - request metrics page
- `POST /admin/reset` - reset demo state in dev mode

## Requirements

- Go 1.20+ (or newer)
- PostgreSQL database
- `sqlc` for query generation (if modifying SQL schema)

## Environment

Set up the environment variables before running the app:

- `DB_URL` - PostgreSQL connection string
- `JWT_KEY_CHIRP` - secret used to sign JWT tokens
- `PLATFORM` - server platform mode (`dev`, `user`, or `admin`)
- `POLKA_KEY` - webhook API key for Polka integration

## Run Locally

1. Install dependencies:
   ```bash
   go mod tidy
   ```
2. Create or update `.env` with the required variables.
3. Start the server:
   ```bash
   go build -o .out && ./.out
   ```
4. Access the app on `http://localhost:8080/app/`

## Development Notes

- The app uses `http.NewServeMux()` and custom handler functions in `server/`.
- `internal/dbman` holds database access logic and `sqlc` query wrappers.
- Passwords are hashed with `argon2id` and JWTs are generated with `github.com/golang-jwt/jwt/v5`.
- Request metrics are tracked in-memory and returned under `/admin/metrics`.
- Text validation filters a small set of profane words before saving chirps.

## Course Context

This repository is part of the Botdev guided course "HTTP servers with Go." It is designed to teach real-world server fundamentals, authentication flows, database integration, and middleware-style metrics handling using native Go libraries.
