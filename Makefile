## Запуск линтера
lint:
	@echo "Running golang-lint..."
	@golangci-lint run

## Запуск сервиса
run:
	@echo "Starting service..."
	@go run cmd/app/main.go

## Сборка бинарника
build:
	@echo "Building the project..."
	@go build -o usdt-rate-service ./cmd/app

## Запуск тестов
test:
	@echo "Running tests..."
	@go test -v ./tests

## Сборка Docker-образа с приложением
up:
	@echo "Starting docker-compose..."
	@docker-compose up -d --build

## Останоыка контейнеров
down:
	@echo "Stopping docker-compose..."
	@docker-compose down

help:
	@echo "Available commands"
	@echo "  make lint          - Run linters"
	@echo "  make build         - Build the project"
	@echo "  make run           - Run the application"
	@echo "  make up            - Start the project with Docker Compose"
	@echo "  make down          - Stop Docker Compose"
	@echo "  make test          - Run all tests"