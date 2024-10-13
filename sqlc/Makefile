.PHONY: migrate sqlc

generate:
	sqlc generate

migrate-up:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/courses" up

migrate-down:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/courses" up

dc-up:
	docker-compose up -d

dc-down:
	docker-compose down