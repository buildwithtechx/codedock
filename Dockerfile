FROM node:22-alpine AS dashboard-builder
WORKDIR /app/dashboard
COPY dashboard/package*.json ./
RUN npm ci
COPY dashboard/ ./
RUN npm run build

FROM golang:1.25-alpine AS daemon-builder
WORKDIR /src
RUN apk add --no-cache git ca-certificates tzdata
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=dashboard-builder /app/dashboard/dist ./dashboard/dist
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s -X main.codedockVersion=${VERSION}" -o /codedockd ./cmd/codedockd

FROM alpine:3.21 AS production
WORKDIR /var/www/codedock
RUN apk add --no-cache ca-certificates tzdata docker-cli git openssh-client curl
COPY --from=daemon-builder /codedockd /usr/local/bin/codedockd
RUN mkdir -p /var/www/codedock/data
ENV PORT=8080 \
    CODEDOCK_DATA_DIR=/var/www/codedock/data
EXPOSE 8080 80 443
VOLUME ["/var/www/codedock/data"]
ENTRYPOINT ["codedockd"]
