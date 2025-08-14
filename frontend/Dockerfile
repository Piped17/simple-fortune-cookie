# Frontend build stage
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

# Build binary
RUN go build -o frontend .

# Runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/frontend /app/frontend
COPY --from=builder /app/static /app/static
COPY --from=builder /app/templates /app/templates

# Set environment variables for backend connection
ENV BACKEND_DNS=backend
ENV BACKEND_PORT=9000

CMD ["./frontend"]
