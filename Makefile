.SILENT: 

.PHONY: fmt lint race test run migrate_up migrate_down migrate_status 

include .env.local
export

fmt:  go fmt ./...

lint: fmt 
	go vet ./...

race: lint 
	go test -v -race ./...

test: race 
	go  test -v -cover ./...

run_source: test 
	go run -v cmd/finances_api/main.go

run: dc_down
	docker compose up -d 

dc_down: 
	docker compose down

migrate_up: 
	goose -dir ./migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=finances_api" up

migrate_down: 
	goose -dir ./migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=finances_api" down

migrate_status: 
	goose -dir ./migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=finances_api" status


.DEFAULT_GOAL := run
