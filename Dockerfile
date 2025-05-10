FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main ./app/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/app/templates ./templates

EXPOSE 8080

CMD ["./main"]
