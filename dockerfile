# ---- Build stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_PATH

RUN go build -o app ${SERVICE_PATH}

# ---- Runtime stage ----
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]