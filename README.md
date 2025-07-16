<img src="./client/src/static/img/monsoon_logo.png"  width="400">

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

Regenerate Swagger docs
```
swag init --dir cmd/api,controller --pd 
```

Run tests
```
go test ./tests/*** -v
```

Generate mocks
```
mockgen -source=[file] -destination=mocks/[destination file] -package=mocks
```

## To Do
- Postgres DB setup for authentication and data store
- Use separate databases for storing application data (users, rooms) and messages
- Implement message bucketing and partitons
- Set up JWT for user auth state
- Write unit tests for DB functions and room states
