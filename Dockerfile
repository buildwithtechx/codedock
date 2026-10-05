FROM node:22-alpine AS dashboard-builder
WORKDIR /app

COPY package*.json ./
COPY apps/ ./apps/
COPY packages/ ./packages/
COPY tsconfig.base.json ./
RUN npm ci

RUN npm run build:dashboard

FROM golang:1.26-alpine AS daemon-builder
WORKDIR /src

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=dashboard-builder /app/apps/dashboard/dist ./apps/dashboard/dist

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s -X main.codedockVersion=${VERSION}" -o /codedockd ./cmd/codedockd

FROM alpine:3.21 AS production
WORKDIR /codedock

RUN apk add --no-cache ca-certificates tzdata docker-cli git openssh-client curl

COPY --from=daemon-builder /codedockd /usr/local/bin/codedockd
RUN mkdir -p /codedock/data

ENV PORT=8080 \
    CODEDOCK_DATA_DIR=/codedock/data \
    CODEDOCK_CLOUD_MODE=false

EXPOSE 8080

VOLUME ["/codedock/data"]

ENTRYPOINT ["codedockd"]
