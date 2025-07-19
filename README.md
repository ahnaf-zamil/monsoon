<img src="./client/src/static/img/monsoon_logo.png"  width="400" style="margin-bottom: 20px">

[![Server CI](https://github.com/ahnaf-zamil/ws_rt_app/actions/workflows/server-ci.yml/badge.svg)](https://github.com/ahnaf-zamil/ws_rt_app/actions/workflows/server-ci.yml)
![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)

## Technology

```
- Go (Gin-Gonic, Gorilla/Websocket)
- TypeScript (React, Tailwind)
- NATS
- PostgreSQL
```

## Features

- Threadsafe global client state for keeping track of sockets
- In-memory "Rooms" implementation for group chats
- Uses NATS for dispatching new messages to other server instances
- JWT-based authentication

## Usage

### Client

Build CSS

```
yarn build:css
```

Run dev server

```
yarn dev
```

Build

```
yarn build
```

### Server

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

API Docs at `http://localhost:9000/swagger/index.html`

## To Do

- Set up separate DB for persisting messages
- Separate "groups" and "dm" in WebSocket rooms state
- Implement message bucketing and partitons
- Implement getter method for user auth
- Reimplement login route
- Clean up register route controller
- Implement login features on client: Fetching encrypted seed, deriving keys from password, deriving signing and encryption keypair on login.

## License

This project is licensed under the AGPL-3.0 License - see the [LICENSE](./LICENSE) file for details.
