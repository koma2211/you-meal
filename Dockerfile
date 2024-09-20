FROM golang:1.23.1-alpine3.20 AS builder 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

FROM alpine3.13

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/configs .
COPY --from=builder /app/.env . 
COPY --from=builder /app/database .

EXPOSE 8080
