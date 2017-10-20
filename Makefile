GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)

dev:
	go fmt .
	go vet .
	go test .

install: dev
	go install .

lint:
	@if gofmt -l . | egrep -v ^vendor/ | grep .go; then \
	  echo "^- Repo contains improperly formatted go files; run make dev" && exit 1; \
	  else echo "All .go files formatted correctly"; fi
	go vet *.go

test: test-all

test-all:
	@go test -v -race -cover -coverprofile=./coverage.out $(GOPACKAGES)

build:
	GOOS=darwin GOARCH=amd64 go build -o purr_darwin_amd64 .
	GOOS=linux GOARCH=amd64 go build -o purr_linux_amd64 .
	GOOS=windows GOARCH=amd64 go build -o purr_windows_amd64 .
