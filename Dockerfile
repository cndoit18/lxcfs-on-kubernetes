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
FROM golang:1.24.4 AS builder

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update; \
    apt-get install -y --no-install-recommends xz-utils

ARG UPX_VERSION=5.0.1
RUN wget -q "https://github.com/upx/upx/releases/download/v$UPX_VERSION/upx-$UPX_VERSION-$(go env GOARCH)_linux.tar.xz" -O - | \
    tar -xJvf - -C /usr/bin --strip-components=1 "upx-$UPX_VERSION-$(go env GOARCH)_linux/upx"

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.sum ./
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY . .

ARG GOLDFLAGS=
ARG CGO_ENABLED=0
RUN go build -ldflags="${GOLDFLAGS}" -a -o manager cmd/manager/main.go

RUN upx -9 manager

# alpine:3.22.0
FROM alpine@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715 
WORKDIR /
COPY --from=builder /workspace/manager .

ENTRYPOINT ["/manager"]
