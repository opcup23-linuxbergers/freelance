all: vet test
vet: go.sum
	go vet ./...
test:
	go test -v ./tests/
fmt:
	go fmt ./...
build: go.sum
	go build -o out
go.sum: go.mod
	go mod tidy
	go mod download
