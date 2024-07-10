deps:
	@go install golang.org/x/tools/cmd/godoc@latest

docs: deps
	@godoc -http=:6060

tidy:
	@go mod tidy

tests:
	@go test -cover ./...

install:
	@go install
