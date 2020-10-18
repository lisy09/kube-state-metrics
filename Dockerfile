# FROM gcr.io/distroless/static
FROM debian:buster-slim

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

COPY /bin/kube-state-metrics-${TARGETOS}-${TARGETARCH}${TARGETVARIANT} /kube-state-metrics

USER nobody

ENTRYPOINT ["/kube-state-metrics", "--port=8080", "--telemetry-port=8081"]

EXPOSE 8080 8081
