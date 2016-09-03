
install:
	go fmt ./...
	go vet ./...
	go test ./...
	go install .

build:
	GOOS=darwin GOARCH=amd64 go build -o purr_darwin_amd64 .
	GOOS=linux GOARCH=amd64 go build -o purr_linux_amd64 .
	GOOS=windows GOARCH=amd64 go build -o purr_windows_amd64 .
