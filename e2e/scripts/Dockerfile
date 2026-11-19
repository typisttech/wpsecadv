ARG COMPOSER_VERSION

FROM docker.io/library/golang:1.27.4-alpine3.24@sha256:70bd7ddea8b6cbe5d6943cde279d6b67feb10cc924fb9100d1d8757c1575b280 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./e2e/scripts ./e2e/scripts
RUN CGO_ENABLED=0 go test -tags=e2escripts -c -o /app/bin/scripts.test ./e2e/scripts

FROM ghcr.io/composer/docker:$COMPOSER_VERSION
ENV TESTSCRIPT_COMPOSER_CACHE_DIR=/tmp/cache
WORKDIR /app

COPY --from=build /app/bin/scripts.test ./scripts.test
COPY ./e2e/scripts/testdata ./testdata

ENTRYPOINT ["/app/scripts.test"]
