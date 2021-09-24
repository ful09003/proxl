run-tests:
	go test -v ./...

precommit:
	go fmt ./...

build:
	go build -o bin/cards cmd/main.go