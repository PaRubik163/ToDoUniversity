# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOMAXPROCS=1 GOGC=50 GOMEMLIMIT=600MiB CGO_ENABLED=0 \
    go build -p 1 -trimpath -buildvcs=false -ldflags="-s -w" -o /out/todo-bot ./cmd/bot

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S bot \
    && adduser -S -G bot bot
WORKDIR /app
COPY --from=build /out/todo-bot /app/todo-bot
USER bot
EXPOSE 8080
ENTRYPOINT ["/app/todo-bot"]
