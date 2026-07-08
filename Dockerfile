FROM docker.io/library/golang:1.26-alpine3.24@sha256:9097beb5536220f7857bdcb65c1b4b340630dd7a70b85f03d5af29640b06693d AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ARG MOD_TIME
ARG REVISION
RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w -X main.ModTime=${MOD_TIME} -X main.Revision=${REVISION}" \
    -o /app/bin/serve ./cmd/serve
FROM gcr.io/distroless/static-debian13:nonroot@sha256:963fa6c544fe5ce420f1f54fb88b6fb01479f054c8056d0f74cc2c6000df5240
EXPOSE 8080

COPY --from=build /app/bin/serve /serve
ENTRYPOINT ["/serve"]
