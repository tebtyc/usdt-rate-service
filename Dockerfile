FROM golang:1.24-alpine AS build

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o usdt-rate-service ./cmd/app

FROM alpine:latest

COPY --from=build /app/usdt-rate-service /usdt-rate-service

CMD ["/usdt-rate-service"]