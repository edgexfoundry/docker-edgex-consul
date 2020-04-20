###############################################################################
# Copyright 2019 Canonical.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
###############################################################################

.PHONY: docker build clean

DOCKERS=edgex-consul
.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
DOCKER_TAG=$(VERSION)

GIT_SHA=$(shell git rev-parse HEAD)

GO=CGO_ENABLED=0 GO111MODULE=on go
GOFLAGS=-ldflags "-X github.com/edgexfoundry/edgex-go.Version=$(VERSION)"


build: health

.PHONY: health
health:
	$(GO) build $(GOFLAGS) -o $@ ./command/health

docker: $(DOCKERS)

edgex-consul:
	 docker build \
		-f Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t edgexfoundry/docker-edgex-consul:$(GIT_SHA) \
		-t edgexfoundry/docker-edgex-consul:$(DOCKER_TAG) \
		.
