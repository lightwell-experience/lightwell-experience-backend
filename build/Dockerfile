FROM registry.access.redhat.com/ubi9/go-toolset:9.8-1790174511 as builder

WORKDIR /go/src/app
COPY . .

USER 0

RUN make get-deps build


FROM registry.access.redhat.com/ubi9/ubi-minimal:latest as ca-source
FROM registry.access.redhat.com/ubi9/ubi-micro:latest

WORKDIR /

COPY --from=ca-source /etc/pki/tls/certs/ca-bundle.crt /etc/pki/tls/certs/

COPY --from=builder /go/src/app/release/* ./

USER 1001

CMD ["/lightwell-network"]
