###############################################################################
# Copyright 2019 Canonical.
# Copyright 2019 Intel Corporation.
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
# Consul image for EdgeX Foundry

# Docker image for Golang health CLI application
FROM golang:1.13-alpine AS builder

ENV GO111MODULE=on
WORKDIR /go/src/github.com/edgexfoundry/docker-edgex-consul

# The main mirrors are giving us timeout issues on builds periodically.
# So we can try these.

RUN sed -e 's/dl-cdn[.]alpinelinux.org/nl.alpinelinux.org/g' -i~ /etc/apk/repositories

RUN apk update && apk add pkgconfig build-base git

COPY go.mod .

RUN go mod download

COPY . .

RUN make build

# consul upstream is based on alpine
FROM consul:1.7.2

LABEL license='SPDX-License-Identifier: Apache-2.0' \
    copyright='Copyright (c) 2019: Canonical; \
Copyright (c) 2019: Intel Corporation'

# for pg_isready to check when kong-db is ready
RUN apk add postgresql-client jq=1.6-r0 curl=7.64.0-r3

# Make sure the default directories are created for consul
RUN mkdir -p /consul/scripts && mkdir -p /consul/config

# Inject EdgeX config/scripts into the image
COPY scripts /edgex/scripts
COPY config /edgex/config
COPY --from=builder /go/src/github.com/edgexfoundry/docker-edgex-consul/health /edgex/scripts/health

# be sneaky and sneak jq into the scripts dir for now, eventually need a 
# statically compiled go program so we don't have to deal with musl/glibc issues
# but for now everything is alpine and thus everything is musl
RUN cp /usr/bin/jq /consul/scripts/jq
RUN cp /usr/lib/libonig.so* /consul/scripts/

# consul ports
EXPOSE 8500
EXPOSE 8400

# Copy over custom entrypoint (derived from consul entrypoint)
COPY edgex-consul-entrypoint.sh /usr/local/bin/edgex-consul-entrypoint.sh
ENTRYPOINT ["edgex-consul-entrypoint.sh"]

# Override consul defaults so that container runs in production mode by default
CMD [ "agent", "-ui", "-bootstrap", "-server", "-client", "0.0.0.0" ]
