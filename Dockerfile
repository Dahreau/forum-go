# Build stage
FROM golang:1.23.1-alpine3.20 AS build

# Set the working directory inside the container
WORKDIR /app

COPY go.mod go.sum ./

# Cache dependencies (unless change in go.mod)
RUN go mod download

COPY . .

RUN go test --cover -v ./...

RUN go build -v -o forum-go

# Final stage: minimal image for development
FROM alpine:latest

LABEL authors="TODO: Add contributors" \
      description="TODO: Add description, maybe?" \
      license="GNU GPL V3.0-or-later" \
      maintainer="See author label" \
      contact="See author label"

WORKDIR /app

COPY --from=build /app/forum-go /app/forum-go

# Copy the source code for development purposes
COPY --from=build /app/src /app/src

EXPOSE 3000

CMD ["/app/forum-go"]

