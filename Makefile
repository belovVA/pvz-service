MOCKERY=/home/vladimir/go/bin/mockery


run:
	go run cmd/pvz-service/main.go

build:
	go build cmd/pvz-service/main.go

build-up:
	docker compose up -d

test:
	go test ./...

generate_repo_mocks:
	$(MOCKERY) --name=UserRepository --dir=internal/service --output=internal/service/mocks
	$(MOCKERY) --name=PvzRepository --dir=internal/service --output=internal/service/mocks
	$(MOCKERY) --name=ReceptionRepository --dir=internal/service --output=internal/service/mocks
	$(MOCKERY) --name=ProductRepository --dir=internal/service --output=internal/service/mocks

generate_service_mocks:
	$(MOCKERY) --name=AuthService --dir=internal/handler --output=internal/handler/mocks
	$(MOCKERY) --name=PvzService --dir=internal/handler --output=internal/handler/mocks
	$(MOCKERY) --name=ReceptionService --dir=internal/handler --output=internal/handler/mocks
	$(MOCKERY) --name=ProductService --dir=internal/handler --output=internal/handler/mocks
	$(MOCKERY) --name=InfoService --dir=internal/handler --output=internal/handler/mocks
