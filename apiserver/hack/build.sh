#!/usr/bin/env bash

CGO_ENABLED=0 go build -o main ./cmd/...
docker build -t "$IMAGE" -f hack/Dockerfile .
