LOCAL_BIN:=$(CURDIR)/bin

run:
	go run cmd/app/main.go --config testdata/config.yaml

build:
	go build -ldflags="-s -w" -trimpath -buildvcs=false -o .bin/app cmd/app/main.go

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/air-verse/air@latest

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

update-packages:
	go get -u ./...

run-air:
	air --build.cmd "go build -o .bin/air cmd/app/main.go" --build.bin "./.bin/air --config testdata/config.yaml"

lint:
	GOBIN=$(LOCAL_BIN) golangci-lint run ./... --config .golangci.yaml

gosec:
	gosec ./...

gocritic:
	gocritic check ./...

test-coverage:
	go test -short -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out