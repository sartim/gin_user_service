# syntax=docker/dockerfile:1
FROM golang:1.27.0-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/user-service ./cmd

FROM alpine:3.24

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app -H app
WORKDIR /app
COPY --from=build --chown=app:app /out/user-service /app/user-service

USER app
EXPOSE 8000
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8000/health/live || exit 1
ENTRYPOINT ["/app/user-service"]
CMD ["--action=run-server"]
