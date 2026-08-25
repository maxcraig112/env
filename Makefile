.PHONY: generate build test vet fmt tidy

generate:
	go generate ./...

build:
	go build ./...

test:
	go test ./... -cover

vet:
	go vet ./...

fmt:
	gofmt -l -w .

tidy:
	go mod tidy
