start:
	go run cmd/main.go

gen:
	protoc --proto_path=proto --go_out=pb --go-grpc_out=pb proto/*.proto

kafka-clean: ## Удалить все данные Kafka
	docker-compose down -v
	@echo "🧹 Kafka data cleaned!"

docker-up:
	docker compose up

docker-down-v:
	docker compose down -v

logs:
	docker compose logs -f app

docker-stop:
	docker compose stop

docker-start:
	docker compose start

test:
	go test -v ./...

docker-rebuild:
	docker compose build --no-cache && docker compose up