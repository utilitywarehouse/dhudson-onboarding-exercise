FROM broady/cacerts
ARG SERVICE

COPY bin/${SERVICE} /app

ENTRYPOINT ["/app"]
