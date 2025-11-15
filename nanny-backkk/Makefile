.PHONY: help run build test clean deps db-setup

help: ## Показать помощь
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Запустить сервер
	go run cmd/api/main.go

build: ## Собрать бинарник
	go build -o bin/nanny-backend cmd/api/main.go

test: ## Запустить тесты
	go test -v ./...

clean: ## Очистить сгенерированные файлы
	rm -rf bin/
	go clean

deps: ## Установить зависимости
	go mod download
	go mod tidy

db-setup: ## Настроить базу данных (требует psql)
	@echo "Создание базы данных..."
	@psql -U postgres -c "CREATE DATABASE nanny_db;" || true
	@echo "Применение схемы..."
	@psql -U postgres -d nanny_db -f scripts/schema.sql

watch: ## Запустить с автоперезагрузкой (требует air)
	air

install-air: ## Установить air для hot reload
	go install github.com/air-verse/air@latest

lint: ## Запустить линтер
	golangci-lint run

fmt: ## Форматировать код
	go fmt ./...
	gofmt -s -w .

vet: ## Проверить код
	go vet ./...

.DEFAULT_GOAL := help
