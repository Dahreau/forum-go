FROM golang:1.23.1-alpine3.20 AS build

WORKDIR /app

COPY . .

RUN go test --cover -v ./...
RUN go build -v -o forum-go

FROM alpine:latest

LABEL authors=""\ #TODO: Add contributor
      description=""\ #TODO: Add description, maybe?
      licence="GNU GPL V3.0-or-later"\
      maintainer="See author label"\
      contact="See author label"

WORKDIR /app
COPY --from=build /app/forum-go /app/forum-go
COPY --from=build /app/src /app/src

EXPOSE 3000

CMD ["/app/forum-go"]
