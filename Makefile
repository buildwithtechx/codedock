BINARY_NAME = codedockd
BUILD_DIR = bin

.PHONY: help all build build-daemon build-dashboard dev dev-dryrun dev-daemon dev-dashboard clean check fmt test docker-build docker-up docker-down

all: check build

help:
	@echo "Codedock — available commands:"
	@echo ""
	@echo "  make check             Run Go fmt + vet and dashboard Biome check"
	@echo "  make fmt               Run Go fmt + dashboard Biome write"
	@echo "  make test              Run Go tests + dashboard unit tests"
	@echo "  make build             Build dashboard + Go daemon binary"
	@echo "  make dev               Run daemon + dashboard dev server concurrently"
	@echo "  make dev-dryrun        Run daemon (DEPLOY_DRY_RUN=true) + dashboard concurrently"
	@echo "  make dev-daemon        Run Go daemon in dev mode"
	@echo "  make dev-dashboard     Run Vite dashboard dev server"
	@echo "  make clean             Remove build artifacts"
	@echo "  make docker-build      Build Docker image"
	@echo "  make docker-up         Start Docker stack"
	@echo "  make docker-down       Stop Docker stack"

check:
	@echo "🔍 Running Go checks & Biome formatting check..."
	go fmt ./...
	go vet ./...
	cd dashboard && npm run format

fmt:
	@echo "🔍 Formatting Go code and dashboard TS/CSS..."
	go fmt ./...
	cd dashboard && npm run format:fix

format: fmt

test:
	@echo "🧪 Running full test suite..."
	go test ./... -v
	cd dashboard && npm run test

build: build-dashboard build-daemon
	@echo "✅ Build complete! Binaries available in $(BUILD_DIR)/ and GUI embedded from dashboard/dist"

build-daemon:
	@echo "⚙️  Building Go daemon binary ($(BINARY_NAME))..."
	mkdir -p $(BUILD_DIR)
	go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/codedockd

build-dashboard:
	@echo "💻 Building React 19 + Vite Dashboard GUI..."
	cd dashboard && npm run build

dev:
	@echo "🚀 Launching Go backend daemon and Vite dashboard GUI concurrently..."
	@$(MAKE) -j2 dev-daemon dev-dashboard

dev-dryrun:
	@echo "🚀 Launching backend daemon (DEPLOY_DRY_RUN=true) and dashboard GUI concurrently..."
	@DEPLOY_DRY_RUN=true $(MAKE) -j2 dev-daemon dev-dashboard

dev-daemon:
	@echo "🚀 Running Go daemon in dev mode with live reload..."
	go run github.com/air-verse/air@latest -c .air.toml

dev-dashboard:
	@echo "💻 Running Dashboard dev server on port 3000..."
	cd dashboard && npm run dev

docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t ghcr.io/buildwithtechx/codedock:dev .

docker-up:
	@echo "🐳 Starting Codedock via Docker..."
	docker network create codedock-network 2>/dev/null || true
	docker run -d --name codedock-control-plane --restart unless-stopped -p 8080:8080 -p 80:80 -p 443:443 -v codedock_data:/codedock/data -e CODEDOCK_DATA_DIR=/codedock/data -v /var/run/docker.sock:/var/run/docker.sock:ro --network codedock-network ghcr.io/buildwithtechx/codedock:dev

docker-down:
	@echo "🐳 Stopping Codedock Docker stack..."
	docker stop codedock-control-plane && docker rm codedock-control-plane

clean:
	@echo "🧹 Cleaning builds and temporary binaries..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	rm -rf dashboard/dist
