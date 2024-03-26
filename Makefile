none:
	@echo ""
build:
	@go build -o bin/api
seed:
	@go run scripts/seed.go
run: build
	@./bin/api
test:
	@go test -v ./... --count=1