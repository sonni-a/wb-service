FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service ./cmd/main/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o producer ./cmd/producer/producer_main.go


FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache bash curl
RUN curl -o wait-for-it.sh https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && \
    chmod +x wait-for-it.sh

COPY --from=builder /app/service .
COPY --from=builder /app/producer .
COPY --from=builder /app/internal/web ./internal/web
COPY --from=builder /app/docs ./docs

RUN chmod +x ./service ./producer

CMD ["./service"]