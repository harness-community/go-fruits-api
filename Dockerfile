#syntax=docker/dockerfile:1.3-labs

FROM golang:1.18-alpine
ARG TARGETARCH

ENV CGO_ENABLED=0

RUN apk add --update

WORKDIR /build

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux GOARCH=${TARGETARCH} go build -o server server.go

CMD [ "/build/server" ]