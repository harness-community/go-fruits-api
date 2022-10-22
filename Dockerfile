#syntax=docker/dockerfile:1.3-labs
FROM goreleaser/goreleaser as builder

WORKDIR build

COPY . .

RUN goreleaser build --single-target --snapshot --rm-dist --output /bin/server

FROM gcr.io/distroless/base

COPY --from=builder /bin/server /bin/server

EXPOSE 8080

ENTRYPOINT ["server"]