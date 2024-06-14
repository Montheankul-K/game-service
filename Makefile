migrate:
	go run databases/migration/migration.go
	go run databases/migrationV2/migrationV2.go
run:
	go run main.go
up:
	docker-compose up -d
down:
	docker-compose down