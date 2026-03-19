# Stage 1: Build Frontend
FROM node:20-slim AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go Server
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy built frontend assets (needed if embedding)
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
RUN go build -o server ./cmd/server/main.go

# Stage 3: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=backend-builder /app/server .
# Re-copy frontend if not embedded or if needed by server
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

EXPOSE 8080
CMD ["./server"]
