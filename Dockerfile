FROM docker.io/library/golang:1.26.1-alpine AS build
WORKDIR /app

# For stamping the binary with git information
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -buildvcs=true -trimpath -ldflags "-s -w" -o /app/bin/server ./cmd/server

FROM gcr.io/distroless/static-debian13:nonroot
EXPOSE 8080

COPY --from=build /app/bin/server /server
ENTRYPOINT ["/server"]
