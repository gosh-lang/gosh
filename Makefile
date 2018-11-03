all: test

init:
	env GO111MODULE=on go mod vendor -v
	go install -v ./vendor/github.com/dvyukov/go-fuzz/go-fuzz \
					./vendor/github.com/dvyukov/go-fuzz/go-fuzz-build \
					./vendor/golang.org/x/tools/cmd/stringer
	curl https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.10.2

install:
	go generate ./...
	go install -v ./...
	go test -i -v ./...

test: install
	go build -v -tags gofuzz ./...
	go test -v ./scanner
	go test -v ./parser
	go test -v ./interpreter
	go test -v -covermode=count -coverprofile=cover.out ./...

check: install
	go run misc/check_license.go
	env GO111MODULE=off golangci-lint run
