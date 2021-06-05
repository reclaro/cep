ifeq ($(shell uname), Darwin)
GOOS=darwin
else
GOOS=linux
endif

GOARCH=amd64


all: vendor test cep

.PHONY: vendor
vendor:
	go mod vendor -v

.PHONY: cep
cep: main.go parsers/*.go printers/*.go utils/*.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build

.PHONY: test
test:
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	go test -coverprofile=coverage.out ./...

.PHONY: view-coverage
view-coverage:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

.PHONY: docker-build
docker-build:
	GOOS=linux go build
	docker build -t cronparserexpander .
