deps:
	@go install golang.org/x/tools/cmd/godoc@latest

docs: deps
	@godoc -http=:6060

tidy:
	@go mod tidy

tests:
	@go test -cover ./...

build:
	@go build -o bin/net-tools ./cmd/net-tools/

build-tui:
	@go build -o bin/net-tui ./cmd/net-tui/

install:
	@go install ./cmd/net-tools/

install-tui:
	@go install ./cmd/net-tui/

install-local:
	@mkdir -p $(HOME)/bin && cp bin/net-tools $(HOME)/bin/net-tools

install-local-tui:
	@mkdir -p $(HOME)/bin && cp bin/net-tui $(HOME)/bin/net-tui

uninstall:
	@rm -f $(HOME)/bin/net-tools

uninstall-tui:
	@rm -f $(HOME)/bin/net-tui

build-linux:
	@GOOS=linux GOARCH=amd64 go build -o bin/net-tools-linux ./cmd/net-tools/

build-linux-tui:
	@GOOS=linux GOARCH=amd64 go build -o bin/net-tui-linux ./cmd/net-tui/

build-windows:
	@GOOS=windows GOARCH=amd64 go build -o bin/net-tools-windows.exe ./cmd/net-tools/

build-windows-tui:
	@GOOS=windows GOARCH=amd64 go build -o bin/net-tui-windows.exe ./cmd/net-tui/

build-mac:
	@GOOS=darwin GOARCH=amd64 go build -o bin/net-tools-mac ./cmd/net-tools/

build-mac-tui:
	@GOOS=darwin GOARCH=amd64 go build -o bin/net-tui-mac ./cmd/net-tui/

build-arm:
	@GOOS=linux GOARCH=arm64 go build -o bin/net-tools-arm ./cmd/net-tools/

build-arm-tui:
	@GOOS=linux GOARCH=arm64 go build -o bin/net-tui-arm ./cmd/net-tui/

build-darwin-universal:
	@GOARCH=amd64 go build -o bin/net-tools-amd64 ./cmd/net-tools/
	@GOARCH=arm64 go build -o bin/net-tools-arm64 ./cmd/net-tools/
	@lipo -create -output bin/net-tools-universal bin/net-tools-amd64 bin/net-tools-arm64
	@rm bin/net-tools-amd64 bin/net-tools-arm64

build-darwin-universal-tui:
	@GOARCH=amd64 go build -o bin/net-tui-amd64 ./cmd/net-tui/
	@GOARCH=arm64 go build -o bin/net-tui-arm64 ./cmd/net-tui/
	@lipo -create -output bin/net-tui-universal bin/net-tui-amd64 bin/net-tui-arm64
	@rm bin/net-tui-amd64 bin/net-tui-arm64

clean:
	@rm -rf bin/

fmt:
	@gofmt -w .

lint:
	@gofmt -d .
