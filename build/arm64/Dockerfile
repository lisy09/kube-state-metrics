FROM golang:1.11.4
RUN mkdir -p /go/src/k8s.io/kube-state-metrics
WORKDIR /go/src/k8s.io/kube-state-metrics
RUN chmod -R 777 /go
COPY / /go/src/k8s.io/kube-state-metrics
COPY /build/arm64/Makefile /go/src/k8s.io/kube-state-metrics/Makefile
RUN make build

FROM alpine:3.8
COPY --from=0 /go/src/k8s.io/kube-state-metrics/kube-state-metrics /kube-state-metrics
ENTRYPOINT ["/kube-state-metrics", "--port=8080", "--telemetry-port=8081"]
EXPOSE 8080 8081
