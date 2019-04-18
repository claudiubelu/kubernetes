# Copyright 2019 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

OS ?= linux
ARCH ?= amd64
TARGET ?= $(CURDIR)
GOARM = 7
GOLANG_VERSION ?= latest
SRC_DIR = $(notdir $(shell pwd))
export

bin:
	docker run --rm -it -v "${TARGET}:${TARGET}:Z" k8s.gcr.io/kube-cross:v1.10.4-1 \
		/bin/bash -c "\
			mkdir -p /go/src /go/bin && \
			go get -d k8s.io/sample-apiserver && \
			cd /go/src/k8s.io/sample-apiserver && \
			git checkout --track origin/release-1.10 && \
			CGO_ENABLED=0 GOARM=${GOARM} GOOS=${OS} GOARCH=${ARCH} go build -a -installsuffix cgo --ldflags '-w' -o ${TARGET}/sample-apiserver k8s.io/sample-apiserver"

.PHONY: bin
