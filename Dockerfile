# ----------------------------
# Builder Stage
# ----------------------------
FROM golang:1.25-alpine AS builder

WORKDIR /app
RUN apk --no-cache add ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build all 3 services
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o avrora-app ./cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o auth-service ./cmd/auth
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o fileserver-service ./cmd/fileserver


# ----------------------------
# Migrate CLI Stage
# ----------------------------
FROM golang:1.25-alpine AS migrate

WORKDIR /app
RUN apk --no-cache add ca-certificates git gcc musl-dev

# Build migrate with postgres driver
RUN CGO_ENABLED=1 GOOS=linux go install \
    -tags 'postgres' \
    -ldflags="-s -w" \
    github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy migrations
COPY internal/db/migrations ./migrations/

# ----------------------------
# Base runtime image (shared)
# ----------------------------
FROM alpine:latest AS base
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY go.mod ./
COPY internal/db/seed/mocks.sql internal/db/seed/
COPY internal/db/migrations/ internal/db/migrations/
COPY image/ image/


# ----------------------------
# Final: avrora-app
# ----------------------------
FROM base AS avrora-app
COPY --from=builder /app/avrora-app .
EXPOSE 8080
CMD ["./avrora-app"]


# ----------------------------
# Final: auth-service
# ----------------------------
FROM base AS auth-service
COPY --from=builder /app/auth-service .
EXPOSE 50051
CMD ["./auth-service"]


# ----------------------------
# Final: fileserver-service
# ----------------------------
FROM base AS fileserver-service
COPY --from=builder /app/fileserver-service .
EXPOSE 50052
CMD ["./fileserver-service"]