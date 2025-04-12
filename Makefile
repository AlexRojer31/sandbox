run:
	go run cmd/sandbox/main.go -c configs/dev.yaml

build:
	mkdir -p bin
	rm -f bin/sandbox
	go build \
		-ldflags "-X github.com/AlexRojer31/sandbox/cmd/sandbox/version.version=$(shell git describe --tags --always --dirty)" \
		-o bin/sandbox \
		cmd/sandbox/main.go

test:
	go test -cover ./...