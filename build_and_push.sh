#!/usr/bin/env bash

set -ex
set -o pipefail

IMAGE_REPO="lisy09kubesphere"
IMAGE_NAME="kube-state-metrics"
IMAGE_VERSION="v1.9.6"
DOCKER_PLATFORMS="linux/amd64,linux/arm/v7,linux/arm64"

echo 'Build and push ks-apiserver'
docker buildx build \
--file Dockerfile \
--tag $IMAGE_REPO/$IMAGE_NAME:$IMAGE_VERSION \
--platform $DOCKER_PLATFORMS \
--push .
