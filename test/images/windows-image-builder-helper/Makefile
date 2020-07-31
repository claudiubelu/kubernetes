# Copyright 2020 The Kubernetes Authors.
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

REGISTRY ?= gcr.io/kubernetes-e2e-test-images
REMOTE_DOCKER_URL ?=
DOCKER_CERT_BASE_PATH ?=
export

build:
	docker --tlsverify --tlscacert "$(DOCKER_CERT_BASE_PATH)/.docker-1809/ca.pem" \
		--tlscert "$(DOCKER_CERT_BASE_PATH)/.docker-1809/cert.pem" --tlskey "$(DOCKER_CERT_BASE_PATH)/.docker-1809/key.pem" \
		-H "$(REMOTE_DOCKER_URL)" build --no-cache --pull -t "$(REGISTRY)/windows-image-builder-helper:1.1" .

push: build
	docker --tlsverify --tlscacert "$(DOCKER_CERT_BASE_PATH)/.docker-1809/ca.pem" \
		--tlscert "$(DOCKER_CERT_BASE_PATH)/.docker-1809/cert.pem" --tlskey "$(DOCKER_CERT_BASE_PATH)/.docker-1809/key.pem" \
		-H "$(REMOTE_DOCKER_URL)" push  "$(REGISTRY)/windows-image-builder-helper:1.1"

.PHONY: build push
