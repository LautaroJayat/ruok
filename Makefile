include .env
export $(shell sed 's/=.*//' .env)

migrate:
	go run cmd/migrate/main.go

seed:
	go run cmd/seed/main.go

build:
	go build cmd/scheduler/main.go

start-db:
	docker compose -f Dockercompose.dev.yml up --build

stop-db:
	docker compose -f Dockercompose.dev.yml down

test:
	go test -count=1 ./...

run:
	./main