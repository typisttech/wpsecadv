FROM docker.io/library/golang:1.26.1-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && \
    go mod verify

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /app/bin/server ./cmd/server

FROM gcr.io/distroless/static-debian13:nonroot

COPY --from=build /app/bin/server /server
ENTRYPOINT ["/server"]
