#!/bin/sh

go build -o ./bin/server ./cmd/api/main.go
go build -o ./bin/migrate ./cmd/migrate/main.go