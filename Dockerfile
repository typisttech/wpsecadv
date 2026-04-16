FROM docker.io/library/golang:1.26.2-alpine3.23@sha256:27f829349da645e287cb195a9921c106fc224eeebbdc33aeb0f4fca2382befa6 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ARG MOD_TIME
ARG REVISION
RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w -X main.ModTime=${MOD_TIME} -X main.Revision=${REVISION}" \
    -o /app/bin/serve ./cmd/serve
FROM gcr.io/distroless/static-debian13:nonroot@sha256:e3f945647ffb95b5839c07038d64f9811adf17308b9121d8a2b87b6a22a80a39
EXPOSE 8080

COPY --from=build /app/bin/serve /serve
ENTRYPOINT ["/serve"]
