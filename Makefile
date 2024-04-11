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
deploy-auth:
	@aws cloudformation deploy --template-file template.yaml --stack-name auth-stack
delete-auth:
	@aws cloudformation delete-stack --stack-name auth-stack
