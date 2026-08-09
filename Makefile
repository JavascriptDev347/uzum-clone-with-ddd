include .env
export

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f postgres

run:
	go run cmd/api/main.go

test:
	go test ./...
swagger:
	swag init -g cmd/api/main.go -o docs

help:
	@echo "Usage:"
	@echo "  make up      - Start the services"
	@echo "  make down    - Stop the services"
	@echo "  make logs    - View logs"
	@echo "  make run     - Run the application"
	@echo "  make test    - Run tests"
	@echo "  make help    - Show this help message"
