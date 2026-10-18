FROM node:22-alpine AS dashboard-builder
WORKDIR /app

COPY package.json package-lock.json ./
COPY apps/dashboard/package.json ./apps/dashboard/
COPY packages/setup/package.json ./packages/setup/
COPY packages/ui/package.json ./packages/ui/
RUN npm ci --workspace=@codedock/dashboard --include-workspace-root

COPY packages/setup/ ./packages/setup/
COPY packages/ui/ ./packages/ui/
COPY apps/dashboard/ ./apps/dashboard/
COPY tsconfig.base.json ./

RUN npm run build:dashboard

FROM golang:1.26-alpine AS daemon-builder
WORKDIR /src

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY apps/dashboard/*.go ./apps/dashboard/
COPY --from=dashboard-builder /app/apps/dashboard/dist ./apps/dashboard/dist

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s -X main.codedockVersion=${VERSION}" -o /codedockd ./cmd/codedockd

FROM alpine:3.21 AS production
WORKDIR /codedock

RUN apk add --no-cache ca-certificates tzdata docker-cli docker-cli-compose git openssh-client curl

COPY --from=daemon-builder /codedockd /usr/local/bin/codedockd
RUN mkdir -p /codedock/data

ENV PORT=8080 \
    CODEDOCK_DATA_DIR=/codedock/data \
    CODEDOCK_CLOUD_MODE=false

EXPOSE 8080

VOLUME ["/codedock/data"]

ENTRYPOINT ["codedockd"]
