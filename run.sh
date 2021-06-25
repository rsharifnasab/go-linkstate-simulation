#!/usr/bin/env bash 

set -o errexit
set -o nounset

find . -type f -name "*.log" -delete

clear 
(
    cd router
    go build
)


(
    cd manager
    go build 
    ./manager
)

cat ./*.log
