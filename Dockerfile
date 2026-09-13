FROM golang:1.23-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/todo-bot .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S bot \
    && adduser -S -G bot bot
WORKDIR /app
COPY --from=build /out/todo-bot /app/todo-bot
USER bot
EXPOSE 8080
ENTRYPOINT ["/app/todo-bot"]
