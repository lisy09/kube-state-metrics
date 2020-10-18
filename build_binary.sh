#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOTDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

TARGET_BINARY=kube-state-metrics
OUTPUT_DIR=bin
BUILDPATH=.

GOBINARY=${GOBINARY:-go}

BUILD_PLATFORMS=${DOCKER_PLATFORMS:-linux/amd64,linux/arm/v7,linux/arm64}
for BUILD_PLATFORM in $(echo ${BUILD_PLATFORMS}| tr "," "\n")
do
    echo "build binary for $BUILD_PLATFORM"
    BUILD_VARS=(${BUILD_PLATFORM//// })
    BUILD_GOOS=${BUILD_VARS[0]}
    BUILD_GOARCH=${BUILD_VARS[1]}
    if [[ -v BUILD_VARS[2] ]]; then
        BUILD_VARIANT=${BUILD_VARS[2]}
    else
        BUILD_VARIANT=""
    fi

    time GOOS=${BUILD_GOOS} CGO_ENABLED=0 GOARCH=${BUILD_GOARCH} ${GOBINARY} build \
            -o ${OUTPUT_DIR}/${TARGET_BINARY}-${BUILD_GOOS}-${BUILD_GOARCH}${BUILD_VARIANT} \
            ${BUILDPATH}
done