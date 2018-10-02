all: test

init:
	go get -u -v golang.org/x/tools/cmd/stringer \
					github.com/dvyukov/go-fuzz/go-fuzz-build \
					github.com/dvyukov/go-fuzz/go-fuzz \
					github.com/golangci/golangci-lint/cmd/golangci-lint

	go get -t -v ./...

install:
	go generate ./...
	go install -v ./...
	go test -i -v ./...
	git diff --exit-code objects/type_string.go

test: install
	go build -v -tags gofuzz ./...
	go test -v ./scanner
	go test -v ./parser
	go test -v ./interpreter
	go test -v -covermode=count -coverprofile=cover.out ./...

check: install
	go run misc/check_license.go
	golangci-lint run ./...
