FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN MAIN_PKG=$(find . -name main.go | head -n 1 | xargs dirname); if [ -z "$MAIN_PKG" ]; then MAIN_PKG="."; fi; CGO_ENABLED=0 GOOS=linux go build -o /app/server $MAIN_PKG

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server /app/server

EXPOSE 3000

CMD ["/app/server"]
