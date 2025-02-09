.SILENT: 

.PHONY: fmt lint race test run_source run migrate_up migrate_down migrate_status dc_down dc_up

include .env.local
export

fmt:  go fmt ./...

lint: fmt 
	go vet ./...

race: lint 
	go test -v -race ./...

dc_up: dc_down
	docker compose up -d 

dc_down: 
	docker compose down

test: race 
	go  test -v -cover ./...

prepare: test
	echo "prepare stage..."

run: dc_up 
	go run cmd/finance_api/main.go

version:
	go version

migrate_up: 
	goose -dir ./migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=finance_api" up

migrate_down: 
	goose -dir ./migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=finance_api" down

migrate_status: 
	goose -dir ./migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=finance_api" status


.DEFAULT_GOAL := run
