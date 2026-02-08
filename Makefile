.PHONY: help run build docker-up docker-down migrate-up migrate-down test lint clean

help: ## Показать справку
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Запустить приложение
	CONFIG_PATH=config/local.yaml go run cmd/server/main.go

build: ## Собрать бинарник
	go build -o bin/server.exe ./cmd/server

docker-up: ## Запустить PostgreSQL в Docker
	docker-compose up -d

docker-down: ## Остановить Docker контейнеры
	docker-compose down

migrate-up: ## Применить миграции
	goose -dir migrations postgres "postgresql://taskuser:taskpass@localhost:5432/taskdb?sslmode=disable" up

migrate-down: ## Откатить последнюю миграцию
	goose -dir migrations postgres "postgresql://taskuser:taskpass@localhost:5432/taskdb?sslmode=disable" down

migrate-status: ## Показать статус миграций
	goose -dir migrations postgres "postgresql://taskuser:taskpass@localhost:5432/taskdb?sslmode=disable" status

test: ## Запустить тесты
	go test -v ./...

lint: ## Запустить линтер
	golangci-lint run

clean: ## Очистить бинарники
	rm -rf bin/

deps: ## Установить зависимости
	go mod download
	go mod tidy

dev: docker-up migrate-up run ## Полный запуск для разработки

