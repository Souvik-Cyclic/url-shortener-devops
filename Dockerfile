# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener main.go

# Runtime Stage
FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/url-shortener /

EXPOSE 8080

CMD ["/url-shortener"]