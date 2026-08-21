FROM docker.io/library/golang:1.27-alpine3.24@sha256:4c9fe60190a2a3350ddc51de80d0224b8a6698d12bdfc999fee45ea9d6c46dbc AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ARG MOD_TIME
ARG REVISION
RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w -X main.ModTime=${MOD_TIME} -X main.Revision=${REVISION}" \
    -o /app/bin/serve ./cmd/serve
FROM gcr.io/distroless/static-debian13:nonroot@sha256:1c2c046bc09ed40fad370b599a0b1ae7987f55b01e247cf27a7c27cd97e5bbc7
EXPOSE 8080

COPY --from=build /app/bin/serve /serve
ENTRYPOINT ["/serve"]
