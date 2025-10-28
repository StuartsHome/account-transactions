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
mocks:
	mockgen -source=store/store.go -destination mocks/mock_store.go

## Tests.
test: mocks
	go test -v ./...

## Docs.
docs:
	rm -rf docs
	swag init
	swag fmt