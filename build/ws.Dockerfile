########################
# 1. Build binary
########################
FROM golang:1.24.5-alpine AS builder

WORKDIR /src
# Only copy go.mod/sum first to leverage Docker layer cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build just the ws-server command
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/ws-server ./cmd/ws-server

########################
# 2. Minimal runtime image
########################
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /out/ws-server /app/
USER nonroot:nonroot
EXPOSE 8085
ENTRYPOINT ["/app/ws-server"]
