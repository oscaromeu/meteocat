ARG  BUILDER_IMAGE=golang:buster
ARG  DISTROLESS_IMAGE=gcr.io/distroless/base
############################
# STEP 1 build executable binary
############################

# syntax = docker/dockerfile:1-experimental

FROM --platform=${BUILDPLATFORM} ${BUILDER_IMAGE} AS base
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download 

FROM base AS build
ARG TARGETOS
ARG TARGETARCH
# Ensure ca-certficates are up to date
RUN update-ca-certificates
WORKDIR /src
COPY . . 
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/meteocat/example/influxdb .


############################
# STEP 2 build a small image
############################
# using base nonroot image
# user:group is nobody:nobody, uid:gid = 65534:65534
FROM ${DISTROLESS_IMAGE}

# Copy our config file and static executable
COPY config.yml $GOPATH/src/meteocat/example/influxdb/
COPY --from=build /go/bin/meteocat/example/influxdb /go/bin/meteocat/example/influxdb
WORKDIR $GOPATH/src/meteocat/example/influxdb/
ENV LOG_DEBUG=true
# Run main
ENTRYPOINT ["/go/bin/meteocat/example/influxdb"]
