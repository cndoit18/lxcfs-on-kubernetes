# Copyright Yinan Li <cndoit18@outlook.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
FROM --platform=$BUILDPLATFORM golang:1.24.4 AS builder

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update; \
    apt-get install -y --no-install-recommends xz-utils

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

ARG UPX_VERSION=5.0.1
RUN BUILDARCH="${BUILDPLATFORM##*/}"; \
    UPX_URL="https://github.com/upx/upx/releases/download/v${UPX_VERSION}/upx-${UPX_VERSION}-${BUILDARCH}_linux.tar.xz"; \
    if wget -q "${UPX_URL}" -O /tmp/upx.tar.xz; then \
        tar -xJvf /tmp/upx.tar.xz -C /usr/bin --strip-components=1 "upx-${UPX_VERSION}-${BUILDARCH}_linux/upx"; \
        rm -f /tmp/upx.tar.xz; \
    else \
        echo "UPX not available for BUILDARCH=${BUILDARCH}; skipping compression"; \
    fi

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.sum ./
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY . .

ARG GOLDFLAGS=
ARG CGO_ENABLED=0
RUN CGO_ENABLED="${CGO_ENABLED}" GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" \
    go build -ldflags="${GOLDFLAGS}" -a -o manager cmd/manager/main.go

RUN BUILDARCH="${BUILDPLATFORM##*/}"; \
    if command -v upx >/dev/null 2>&1 && [ "${TARGETARCH}" = "${BUILDARCH}" ]; then \
        upx -9 manager; \
    else \
        echo "Skipping UPX (upx missing or cross-compile TARGETARCH=${TARGETARCH} BUILDARCH=${BUILDARCH})"; \
    fi

# alpine:3.22.0
FROM alpine:3.22.0
WORKDIR /
COPY --from=builder /workspace/manager .

ENTRYPOINT ["/manager"]
