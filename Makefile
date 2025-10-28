## Docker.
start:
	docker compose up -d

start-db:
	docker compose up -d db

stop:
	docker compose down --volumes
	rm -rf bin/main

restart: stop start

## Local.
build:
	go build -o bin/main

start-local: build
	./bin/main

## Mocks.
remove-mocks:
	rm -rf mocks/*

mocks: remove-mocks
	mockgen -source=store/store.go -destination mocks/mock_store.go

## Tests.
test: mocks
	go test -v ./...

## Docs.
remove-docs:
	rm -rf docs

docs: remove-docs
	swag init
	swag fmt