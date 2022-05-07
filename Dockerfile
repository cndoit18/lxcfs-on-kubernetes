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
FROM golang:1.17 as builder

ENV UPX_VERSION 3.96
RUN apt update \
    && apt install -y xz-utils \
    && wget -c https://github.com/upx/upx/releases/download/v$UPX_VERSION/upx-$UPX_VERSION-$(go env GOARCH)_$(uname -s | tr '[:upper:]' '[:lower:]').tar.xz -O - | \
    tar xJf - --strip-components 1 upx-$UPX_VERSION-$(go env GOARCH)_$(uname -s | tr '[:upper:]' '[:lower:]')/upx \
    && mv ./upx /usr/bin/upx

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY . .

ARG GOLDFLAGS
RUN CGO_ENABLED=0 go build -ldflags="${GOLDFLAGS}" -a -o manager cmd/manager/main.go \
    && upx -9 manager

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/base@sha256:d8244d4756b5dc43f2c198bf4e37e6f8a017f13fdd7f6f64ec7ac7228d3b191e
WORKDIR /
COPY --from=builder /workspace/manager .

ENTRYPOINT ["/manager"]
