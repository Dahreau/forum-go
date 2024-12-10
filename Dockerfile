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

# Create an empty SQLite database file and populate it
RUN sqlite3 /app/db.sqlite ".databases" && \
    sqlite3 /app/db.sqlite < /app/query.sql

# Final stage: minimal image for running the application
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache sqlite

# Copy the binary and assets from the build stage
COPY --from=build /app/forum-go .
COPY --from=build /app/assets /app/assets
COPY --from=build /app/query.sql .
COPY --from=build /app/key.pem .
COPY --from=build /app/cert.pem .
COPY --from=build /app/server.key .
COPY --from=build /app/server.crt .
COPY --from=build /app/.env .
# Copy the prebuilt SQLite database
COPY --from=build /app/db.sqlite /app/db-init.sqlite

# Ensure the binary is executable
RUN chmod +x /app/forum-go

EXPOSE 8080

# Ensure the database is copied to the volume at runtime
CMD ["/bin/sh", "-c", "if [ ! -f /app/db.sqlite ]; then cp /app/db-init.sqlite /app/db.sqlite; fi && ./forum-go"]
