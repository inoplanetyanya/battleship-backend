.PHONY:
.SILENT:

run:
	go run ./cmd/app/main.go

rundb:
	docker run --name='battleship-db' -e POSTGRES_PASSWORD='qwerty' -p 5433:5432 --rm -d postgres