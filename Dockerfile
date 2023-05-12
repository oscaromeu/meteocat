FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.15 as builder
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ARG Version
ARG GitCommit

ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /go/src/github.com/oscaromeu/meteocat/examples/influx
COPY .  .

RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
  go build -ldflags \
  "-s -w -X 'main.Version=${Version}' -X 'main.GitCommit=${GitCommit}'" \
  -a -installsuffix cgo -o /usr/bin/meteocat .

FROM --platform=${BUILDPLATFORM:-linux/amd64} gcr.io/distroless/static:nonroot
# Copy our config file and static executable
# Run main
WORKDIR /
COPY config.yml .
COPY --from=builder /usr/bin/meteocat /
USER nonroot:nonroot
CMD ["/meteocat"," -c config.yml"]
