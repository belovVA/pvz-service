# Убираем моки из покрытия
PKGS=$(shell go list ./... | grep -vE '/(test)')
COVERPKG=$(shell go list ./... | grep -vE '/(mocks|test)' | paste -sd, -)

.PHONY: build-up test cover

build-up:
	docker compose up -d

test:
	go clean -testcache
	go test -covermode=atomic -coverpkg=$(COVERPKG) -coverprofile=coverage.out $(PKGS)

integrate_test:
	# Run the Go service in the background
	go test test/integration_test.go

cover:
	go tool cover -func=coverage.out
