FROM golang:1.24.2 AS builder         
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./app

FROM debian:bookworm-slim            
WORKDIR /app
COPY --from=builder /src/main /app/
COPY --from=builder /src/app/templates /app/templates
EXPOSE 8080
ENTRYPOINT ["/app/main"]
