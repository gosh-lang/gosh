all: test

init:
	go get -u golang.org/x/tools/cmd/stringer \
				github.com/dvyukov/go-fuzz/... \
				github.com/golangci/golangci-lint/cmd/golangci-lint

	go test -i -v ./...

install:
	go generate ./...
	go install -v ./...

test: install
	go build -v -tags gofuzz ./...
	go test -v ./scanner
	go test -v ./parser
	go test -v ./interpreter
	go test -v -covermode=count -coverprofile=cover.out ./...

check: install
	go run misc/check_license.go
	golangci-lint run ./...
