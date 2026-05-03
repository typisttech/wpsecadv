ARG COMPOSER_VERSION

FROM docker.io/library/golang:1.26.2-alpine3.23@sha256:f85330846cde1e57ca9ec309382da3b8e6ae3ab943d2739500e08c86393a21b1 AS build
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
