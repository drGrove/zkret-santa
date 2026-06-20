# syntax=docker/dockerfile:1
FROM --platform=amd64 stagex/pallet-go@sha256:35d936e6fcbd73efc2dc3471a015b1ac5a5bb6d2c0eae6e84550f6352857457d AS pallet-go
FROM --platform=amd64 stagex/core-musl@sha256:42fd2ed18ac4d6336f08f1b6345e0bb7bc5570450a40e7255b6da79868bd2d6b AS musl

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
