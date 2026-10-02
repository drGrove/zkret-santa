# syntax=docker/dockerfile:1
FROM --platform=amd64 stagex/pallet-go@sha256:186bcc2733968cb763b787db230252746c76dc1119529e1b1759b99f2f368fb0 AS pallet-go
FROM --platform=amd64 stagex/core-musl@sha256:94f77e4b1841987a42efce52aac54dbf230c152cb55c62da1967189a2705df89 AS musl

FROM pallet-go AS build
ARG TARGETOS
ARG TARGETARCH
ARG BUILD_TAGS=""
ARG BINARY_NAME=zkret-santa

ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}
ENV EXTRA_ARGS=""

ADD . /zkret-santa
WORKDIR /zkret-santa
RUN go mod download
RUN --network=none <<-EOF
	set -eu
	if [ -n "${BUILD_TAGS}" ]; then
    EXTRA_ARGS="-tags ${BUILD_TAGS}"
	fi
  go build \
    ${EXTRA_ARGS} \
    -trimpath \
    -v \
    -mod=readonly \
    -o ${BINARY_NAME} \
    ./cmd/zkret-santa
	install -Dm0755 -t /rootfs-${TARGETOS}-${TARGETARCH}/usr/bin/ ${BINARY_NAME}
  install -Dm0644 -t /rootfs-${TARGETOS}-${TARGETARCH}/usr/share/licenses/ LICENSE
  install -Dm0644 -t /rootfs-${TARGETOS}-${TARGETARCH}/usr/share/licenses/ COPYRIGHT
EOF

FROM scratch AS package-secret-santa
ARG TARGETOS
ARG TARGETARCH
ARG BINARY_NAME=secret-santa
COPY --from=build /rootfs-${TARGETOS}-${TARGETARCH}/ /
ENTRYPOINT ["/usr/bin/${BINARY_NAME}"]

FROM scratch AS package-zkret-santa
COPY --from=musl / /
ARG TARGETOS
ARG TARGETARCH
ARG BINARY_NAME=zkret-santa
COPY --from=build /rootfs-${TARGETOS}-${TARGETARCH}/ /
ENTRYPOINT ["/usr/bin/${BINARY_NAME}}"]
