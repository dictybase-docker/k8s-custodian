FROM golang:1.16.2-buster AS builder
ENV GOPROXY https://proxy.golang.org
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN apt-get -qq update \
    && apt-get -yqq install upx
RUN mkdir -p /k8s-custodian
WORKDIR /k8s-custodian
COPY go.mod ./
COPY go.sum ./
RUN go mod download
ADD cmd cmd
ADD internal internal
RUN go build \
    -a \
    -ldflags "-s -w -extldflags '-static'" \
    -installsuffix cgo \
    -tags netgo \
    -o /bin/app \
    cmd/k8s-custodian/main.go
RUN strip /bin/app \
    && upx -q -9 /bin/app

FROM gcr.io/distroless/static
ARG BUILD_DATE
ARG VCS_REF
ARG VERSION
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.version=$VERSION \
      org.label-schema.name="k8s-custodian" \
      org.label-schema.description="cli to manage various repetitive tasks in dictybase kubernetes cluster" \
      org.label-schema.vcs-url="https://github.com/dictybase-docker/k8s-custodian" \
      org.label-schema.schema-version="1.0"
COPY --from=builder /bin/app /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/app"]
