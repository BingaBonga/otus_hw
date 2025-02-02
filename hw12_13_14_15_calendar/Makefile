BIN := "./bin/"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./calendar_scheduler
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./calendar_storer

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

up-build:
	docker-compose up -d --build

up:
	docker-compose up

ps:
	docker-compose ps -a

down:
	docker-compose down

integration-test:
	set -e ;\
	docker-compose -f docker-compose.test.yml up --build -d ;\
	docker-compose ps -a ;\
	test_status_code=0;\
	sleep 5;\
	docker-compose -f docker-compose.test.yml run integration_tests go test ./integration-test/... || test_status_code=$$?;\
	docker-compose -f docker-compose.test.yml down;\
	exit $$test_status_code;\

compose:
	docker-compose -f docker-compose.yml up --build -d ;\
	docker-compose ps -a ;\

prune:
	docker system prune -a
