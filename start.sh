#!/bin/bash

# Build the application
go build -o bin/app ./cmd/main.go

# Start the MongoDB Docker container
docker-compose up -d

# Start the application
docker-compose run app