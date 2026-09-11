FROM golang:1.26.7-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

RUN mkdir -p /app/pb_data

COPY --from=builder /app/server /app/server

EXPOSE 8080

CMD ["/app/server", "serve", "--http=0.0.0.0:8080", "--dir=/app/pb_data"]
