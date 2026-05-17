# ── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install templ CLI for code generation
RUN go install github.com/a-h/templ/cmd/templ@v0.2.793

# Cache dependencies separately from source
COPY go.mod go.sum ./
RUN go mod download

# Copy full source
COPY . .

# Generate *_templ.go files from .templ files
RUN templ generate

# Build the binary. CGO_ENABLED=0 keeps it static so the runtime image
# does not need glibc. -ldflags trims debug info and symbol table.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /theology-tracker \
    ./cmd/server

# ── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM alpine:3.19

# ca-certificates is needed if the app ever makes outbound TLS calls.
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /theology-tracker .
COPY --from=builder /app/static ./static

# The data directory is the bind mount target — it must exist in the image
# so Docker can mount into it even if the host path doesn't exist yet.
RUN mkdir -p /app/data

EXPOSE 3000

CMD ["./theology-tracker"]
