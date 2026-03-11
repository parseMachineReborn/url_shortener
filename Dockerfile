FROM golang:1.25.0-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/server/main.go
RUN go build -o migrator ./cmd/migrator/main.go

FROM alpine
WORKDIR /app

COPY --from=builder /app/app ./
COPY --from=builder /app/migrator ./
COPY --from=builder /app/migrations ./migrations

CMD [ "./app" ]