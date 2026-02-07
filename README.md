SchemaWatch

SchemaWatch is a backend service that monitors PostgreSQL table schemas over time, detects schema drift, classifies the impact of changes, and notifies consumers before breaking changes cause production issues.

The system is designed to run as a long lived background service and periodically checks registered tables for schema changes.

Features

Register PostgreSQL tables for monitoring

Capture and version schema snapshots

Detect schema drift between versions

Classify changes as Safe, Risky, or Breaking

Store drift events for audit and history

Trigger webhook notifications on breaking changes

API key based authentication

One time API key generation utility

Dockerized setup for easy local or cloud deployment

Tech Stack

Language: Go

Database: PostgreSQL

Migrations: Goose

Web Framework: Gin

Scheduler: Internal ticker based scheduler

Containerization: Docker

Cloud Database: Neon PostgreSQL

How It Works

A user registers a PostgreSQL table.

The service stores the initial schema snapshot.

A background scheduler periodically fetches the latest schema.

The new schema is compared with the previous version.

Drift is detected and classified.

Drift events are stored.

If a breaking change is detected, a webhook is triggered.

API Authentication

All APIs except the health check require an API key.

The API key must be sent as a header:

X-API-Key: your_api_key_here

API Endpoints

Health check
GET /health

Register a table
POST /v1/schemas

List registered tables
GET /v1/schemas

Get latest schema snapshot
GET /v1/schemas/{schema_id}/latest

Get schema version history
GET /v1/schemas/{schema_id}/versions

Get schema drift events
GET /v1/schemas/{schema_id}/drifts

API Key Generation

A one time utility is provided to generate API keys.

This command generates a raw API key and its hashed value. The hashed value should be stored in the database.

Run:

```
go run ./cmd/key
```

Copy the generated hashed key and insert it into the api_keys table.

Running with Go

To run the application directly using Go, set the DATABASE_URL environment variable and start the server.

```
go run ./cmd
```

The service will start on port 8080.

Running with Docker

The service is fully dockerized and designed to run with an external PostgreSQL database.

Environment variables are used for configuration.

DATABASE_URL must be provided and should point to a PostgreSQL database.

Example environment file:

```
DATABASE_URL=postgres://user:password@host:5432/dbname?sslmode=require
```

To start the service using Docker:

```
docker-compose up --build
```

Database Migrations

Database schema migrations are managed using Goose.

Migrations should be applied before running the service.

Background Scheduler

SchemaWatch runs an internal scheduler that periodically:

Fetches current table schemas

Creates new schema versions

Detects and stores drift events

Triggers webhook notifications on breaking changes

Because of this background processing, the service is intended to run continuously and is not suitable for free hosting platforms that suspend inactive services.

Webhook Notifications

Webhook notifications are triggered only for breaking schema changes.

The webhook payload includes:

Schema ID

Change summary

Impact level

Version range
