# Real-Time Websocket Chat App

[![Go CI](https://github.com/ahnaf-zamil/ws_rt_app/actions/workflows/ci.yml/badge.svg)](https://github.com/ahnaf-zamil/ws_rt_app/actions/workflows/ci.yml)

## Technology
```
- Go (Gin-Gonic, Gorilla/Websocket)
- NATS
- PostgreSQL
```

## Features
- Threadsafe global client state for keeping track of sockets
- In-memory "Rooms" implementation for group chats
- Uses NATS for dispatching new messages to other server instances

## Usage

Start server
```
go run cmd/api/main.go
```

Generate DB schema (will extend for migration later)
```
go run cmd/migrate/main.go
```

Run tests
```
go test ./tests/*** -v
```

## To Do
- Integrate Kafka for streaming new messages to batch processors (This is to prevent bombarding the DB with new messages)
- Create background workers or worker services to batch write new messages to PostgreSQL database
- Postgres DB setup for authentication and data store
- Use separate databases for storing application data (users, rooms) and messages
- Implement message bucketing and partitons
- Set up JWT for user auth state
- Write unit tests for DB functions and room states
