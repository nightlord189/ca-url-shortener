run:
	docker-compose up -d --force-recreate --build

test:
	go test ./...

lint:
	gofumpt -l -w .
	golangci-lint run --timeout=5m

generate:
	swag init --dir ./cmd/app --parseDependency --parseInternal
	go generate ./...