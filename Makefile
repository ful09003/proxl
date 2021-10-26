run-tests:
	go test -v ./...

run-e2e: build
	$(MAKE) -C ./e2e run-e2e
	$(MAKE) -C ./e2e docker-compose-cleanup

precommit:
	go fmt ./...

build:
	go build -o bin/cards cmd/main.go