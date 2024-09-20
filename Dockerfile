FROM golang:alpine AS stage1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

# stage2
FROM alpine
WORKDIR /

COPY --from=stage1 app/main /main
COPY --from=stage1 app/.env /.env

EXPOSE 8099

CMD [ "/main" ]