# Build stage
FROM golang:1.23.1-alpine3.20 AS build

WORKDIR /app

# Install required libraries for CGO and SQLite
RUN apk add --no-cache sqlite gcc musl-dev

# Enable CGO for go-sqlite3
ENV CGO_ENABLED=1

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN go build -v -o forum-go ./cmd/api

# Create an empty SQLite database file
RUN sqlite3 /app/db.sqlite ""

# Final stage: minimal image for running the application
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache sqlite gcc musl-dev

# Copy the binary from the build stage
COPY --from=build /app/forum-go .
COPY --from=build /app/assets /app/assets

# Copy the SQLite database file
COPY --from=build /app/db.sqlite .

# Copy the SQL query file
COPY --from=build /app/query.sql .

# Ensure the binary is executable
RUN chmod +x /app/forum-go

EXPOSE 8080

CMD ["./forum-go"]
