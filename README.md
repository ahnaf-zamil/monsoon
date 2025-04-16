# Real-time websocket chat app 

## Stuff Used
```
- Go (Gin-Gonic, Gorilla/Websocket)
- NATS
```

## Features
- Threadsafe global client state for keeping track of sockets
- In-memory "Rooms" implementation for group chats
- Uses NATS for dispatching new messages to other server instances

## To Do
- Integrate Kafka for streaming new messages to batch processors (This is to prevent bombarding the DB with new messages)
- Create background workers or worker services to batch write new messages to PostgreSQL database
- Postgres DB setup for authentication and data store
