default:
    just --list
build:
    #!/usr/bin/env bash
    set -euxo
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o test_cert .
    docker build -t cert-test -f Dockerfile-1.24 .
