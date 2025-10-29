start:
	go run main.go

start-multiple:
	@echo "Starting servers..."
	PORT=8081 go run main.go &
	PORT=8082 go run main.go &
	PORT=8083 go run main.go &
	wait

stop-multiple:
	@echo "Stopping servers..."
	@lsof -ti tcp:8081 | xargs kill -9 || true
	@lsof -ti tcp:8082 | xargs kill -9 || true
	@lsof -ti tcp:8083 | xargs kill -9 || true

start-db:
	docker-compose up -d

stop-db:
	docker-compose down

lint: lint/install lint/run

lint/install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v2.5.0

lint/run:
	bin/golangci-lint run --config .golangci.yaml