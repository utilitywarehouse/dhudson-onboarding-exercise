FROM alpine:3.8
ARG SERVICE

RUN apk add --no-cache ca-certificates

COPY bin/${SERVICE} /app

ENTRYPOINT ["/app"]
