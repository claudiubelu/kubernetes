#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

REMOTE_DOCKER_URL_1903=${1:-""}
REMOTE_DOCKER_URL_1809=${2:-""}
REGISTRY=${3:-"gcr.io/kubernetes-e2e-test-images"}

REMOTE_DOCKER_URL_1903=$REMOTE_DOCKER_URL_1903 REMOTE_DOCKER_URL_1809=$REMOTE_DOCKER_URL_1809 REGISTRY=$REGISTRY make all-push WHAT=busybox
REMOTE_DOCKER_URL_1903=$REMOTE_DOCKER_URL_1903 REMOTE_DOCKER_URL_1809=$REMOTE_DOCKER_URL_1809 REGISTRY=$REGISTRY make all-push WHAT=mounttest
REMOTE_DOCKER_URL_1903=$REMOTE_DOCKER_URL_1903 REMOTE_DOCKER_URL_1809=$REMOTE_DOCKER_URL_1809 REGISTRY=$REGISTRY make all-push WHAT=test-webserver

all_images=`ls -d */ | grep -v "volume\|echoserver\|jessie\|node-perf\|pets\|cuda-vector-add\|busybox\|mounttest\|test-webserver" | sed "s|/||g"`
all_images=`ls -d */ | grep -v "volume\|node-perf\|pets\|cuda-vector-add\|busybox\|mounttest\|test-webserver" | sed "s|/||g"`
all_images=`ls -d */ | grep -v "volume\|metadata-concealment\|regression-issue-74839\|node-perf\|pets\|cuda-vector-add\|busybox\|mounttest\|test-webserver" | sed "s|/||g"`

for image in $all_images; do
    REMOTE_DOCKER_URL_1903=$REMOTE_DOCKER_URL_1903 REMOTE_DOCKER_URL_1809=$REMOTE_DOCKER_URL_1809 REGISTRY=$REGISTRY make all-push WHAT=$image
done
