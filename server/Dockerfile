FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/event-planner ./cmd/event-planner

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/bin/event-planner ./bin/event-planner

EXPOSE ${SERVER_PORT}

CMD ["./bin/event-planner"]