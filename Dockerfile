ARG SERVICE

# Build stage
FROM golang:latest as build-env
ARG SERVICE
ARG GO=CGO_ENABLED=0 go

ADD . /go/src/github.com/utilitywarehouse/${SERVICE}
RUN cd /go/src/github.com/heedson/${SERVICE} && ${GO} build -o /${SERVICE}

# Production stage
FROM broady/cacerts
ARG SERVICE
COPY --from=build-env /${SERVICE} /app

ENTRYPOINT ["/app"]