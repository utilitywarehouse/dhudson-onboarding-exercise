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

test:
	$(GO) test ./...

download:
	${GO} mod download

ci-docker-auth:		## Docker login for registry access
	@echo "Logging in to $(DOCKER_REGISTRY) as $(DOCKER_ID)"
	@docker login -u $(DOCKER_ID) -p $(UW_DOCKER_PASS) $(DOCKER_REGISTRY)

ci-docker-build: ci-docker-auth
	docker build -t $(DOCKER_REPOSITORY):$(GIT_HASH) . --build-arg SERVICE=$(SERVICE)
	docker tag $(DOCKER_REPOSITORY):$(GIT_HASH) $(DOCKER_REPOSITORY):latest

ci-docker-push: ci-docker-build
	docker push $(DOCKER_REPOSITORY)

ci-k8s-deploy:
	@echo "Patching k8s deployment"
	@curl -X PATCH -k -d '{"spec":{"template":{"spec":{"containers":[{"name":"$(DOCKER_CONTAINER_NAME)","image":"$(DOCKER_REPOSITORY):$(GIT_HASH)"}]}}}}' -H "Content-Type: application/strategic-merge-patch+json" -H "Authorization: Bearer $(K8S_DEV_TOKEN)" https://elb.master.k8s.dev.uw.systems/apis/apps/v1/namespaces/$(NAMESPACE)/deployments/$(DOCKER_CONTAINER_NAME)
