run:
	go run cmd/pvz-service/main.go

build:
	go build cmd/pvz-service/main.go

build-up:
	docker compose up -d
