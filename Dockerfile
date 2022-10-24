# syntax=docker/dockerfile:1.4

FROM gcr.io/distroless/base

ARG TARGETARCH

COPY server_linux_${TARGETARCH}/server /bin/server

EXPOSE 8080

ENTRYPOINT ["server"]