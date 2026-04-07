.PHONY: dev dev-backend dev-frontend build build-backend build-frontend clean install

# ── Development ──────────────────────────────────────────────

# Run both backend and frontend concurrently
dev:
	@echo "Starting backend (localhost:8080) and frontend (localhost:3000)..."
	@make -j2 dev-backend dev-frontend

dev-backend:
	go run ./cmd/server

dev-frontend:
	cd web && npm run dev

# ── Build ────────────────────────────────────────────────────

build: build-backend build-frontend

build-backend:
	go build -o bin/server ./cmd/server

build-frontend:
	cd web && npm run build

# ── Setup ────────────────────────────────────────────────────

install:
	go mod download
	cd web && npm install

# ── Clean ────────────────────────────────────────────────────

clean:
	rm -rf bin/
	rm -rf web/.next/
