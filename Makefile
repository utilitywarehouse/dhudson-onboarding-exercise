SERVICE=dhudson-onboarding-exercise

DOCKER_ID=onboarding
DOCKER_REGISTRY=registry.uw.systems
DOCKER_CONTAINER_NAME=$(SERVICE)
NAMESPACE=onboarding
DOCKER_REPOSITORY=$(DOCKER_REGISTRY)/$(NAMESPACE)/$(DOCKER_CONTAINER_NAME)

BUILDENV :=
BUILDENV += CGO_ENABLED=0 GO111MODULE=on GOPROXY=direct GOPRIVATE=github.com/utilitywarehouse/*
GO := ${BUILDENV} go
GIT_HASH := $(CIRCLE_SHA1)
ifeq ($(GIT_HASH),)
  GIT_HASH := $(shell git rev-parse HEAD)
endif

lint:			## Run linting
	@which golangci-lint || ${GO} get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.21.0
	${BUILDENV} $(GOPATH)/bin/golangci-lint run ./...

clean:			## Removes compiled artifacts
	rm -rf bin/

build: clean
	$(GO) build -o bin/$(SERVICE) ./


.PHONY: test
test: mock
	$(GO) test ./...


ci-docker-auth:		## Docker login for registry access
	@echo "Logging in to $(DOCKER_REGISTRY) as $(DOCKER_ID)"
	@docker login -u $(DOCKER_ID) -p $(DOCKER_PASSWORD) $(DOCKER_REGISTRY)

ci-docker-build: ci-docker-auth
	docker build -t $(DOCKER_REPOSITORY):$(GIT_HASH) . --build-arg SERVICE=$(SERVICE) --build-arg GO=$(GO)
	docker tag $(DOCKER_REPOSITORY):$(GIT_HASH) $(DOCKER_REPOSITORY):latest

ci-docker-push: ci-docker-build
	docker push $(DOCKER_REPOSITORY)

docker-build: docker-auth
	docker build . -t $(DOCKER_REPOSITORY):$(GIT_HASH) -t $(DOCKER_REPOSITORY):latest

docker-auth:
	@echo "Logging in to $(DOCKER_REGISTRY)"
	@docker login -u $(DOCKER_ID) -p $(UW_DOCKER_PASS) $(DOCKER_REGISTRY)

docker-push:
	docker push $(DOCKER_REPOSITORY):$(GIT_HASH)
	docker push $(DOCKER_REPOSITORY):latest

download:
	${GO} mod download