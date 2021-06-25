#!/usr/bin/env bash 

set -o errexit
set -o nounset

(
    cd router
    go build
)


(
    cd manager
    go build 
    ./manager
)
