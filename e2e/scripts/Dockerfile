ARG COMPOSER_VERSION

FROM docker.io/library/golang:1.26-alpine3.24@sha256:d29e660cc75a5b6b1334e03c5c81ccf9bc0884a002c6000dbf0fb96034814478 AS build
WORKDIR /app

ENV GOEXPERIMENT=jsonv2

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
