#!/usr/bin/env bash 

set -o errexit
set -o nounset

find . -type f -name "*.log" -delete

clear 
(
    cd router
    go build -race
)


(
    cd manager
    go run -race . || true
    #./manager || true
)

cat ./*.log
