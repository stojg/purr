
dev:
	go fmt .
	go vet .
	go test .

install:
	go fmt .
	go vet .
	go test .
	go install .

lint:
	@if gofmt -l . | egrep -v ^vendor/ | grep .go; then \
	  echo "^- Repo contains improperly formatted go files; run make dev" && exit 1; \
	  else echo "All .go files formatted correctly"; fi
	go vet *.go

ci: lint
	go get
	go test .

build:
	GOOS=darwin GOARCH=amd64 go build -o purr_darwin_amd64 .
	GOOS=linux GOARCH=amd64 go build -o purr_linux_amd64 .
	GOOS=windows GOARCH=amd64 go build -o purr_windows_amd64 .
