.PHONY: build

LINTER_VERSION=v1.31.0

default: build

## lint: analyse the source code with the configuration in .golangci.yml
lint:
	golangci-lint run --timeout=5m

unit-tests:
	go fmt ./...
	go test -shuffle=on --tags=unit ./...

## test-race: run tests with race detection
race-condition-tests:
	go test -v -race ./...

integration-tests:
	go test -count=1 --tags=integration ./...

acceptance-tests:
	go test -count=1 --tags=acceptance ./...

docker-tests:
	docker-compose down
	docker-compose build
	docker-compose run --rm unit-tests
	docker-compose run --rm integration-tests
	docker-compose run --rm acceptance-tests
	docker-compose down

build: lint unit-tests race-condition-tests integration-tests acceptance-tests

mod:
	go mod vendor -v

tidy:
	go mod tidy -v

down:
	docker-compose down

up:
	docker-compose up
