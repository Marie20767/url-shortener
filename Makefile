GIT_HASH := $(shell git rev-parse --short HEAD)
DOCKER_USER := marie20767
IMAGE_NAME := url-shortener-server

dep:
	go mod download

start:
	go run main.go

start-multiple:
	@echo "Starting servers..."
	SERVER_PORT=8081 go run main.go &
	SERVER_PORT=8082 go run main.go &
	SERVER_PORT=8083 go run main.go &
	wait

stop-multiple:
	@echo "Stopping servers..."
	@lsof -ti tcp:8081 | xargs kill -9 || true
	@lsof -ti tcp:8082 | xargs kill -9 || true
	@lsof -ti tcp:8083 | xargs kill -9 || true

start-dbs:
	docker compose up -d postgres url-cache key-cache

stop-dbs:
	docker compose down

test:
	go test ./...

lint: lint-install lint-run

lint-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v2.5.0

lint-run:
	bin/golangci-lint run --config .golangci.yaml

lint-fix:
	bin/golangci-lint run --config .golangci.yaml --fix

build:
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -o url-shortener-server .

docker/build-and-push:
	DOCKER_BUILDKIT=1 docker buildx build \
	--platform linux/amd64 \
	--push \
	-t $(DOCKER_USER)/$(IMAGE_NAME):latest \
	-t $(DOCKER_USER)/$(IMAGE_NAME):$(GIT_HASH) .