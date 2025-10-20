default:
    just --list
build-scratch:
    #!/usr/bin/env bash
    set -euxo
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o test_cert .
    docker build -t cert-test -f Dockerfile-1.24 .

build-two-step:
    #!/usr/bin/env bash
    set -euxo
    docker build -t test-cert -f Dockerfile .
