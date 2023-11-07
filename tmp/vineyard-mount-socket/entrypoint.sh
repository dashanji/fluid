#!/bin/sh
set -ex

mkdir -p /runtime-mnt/vineyard/default/vineyard/vineyard-fuse
while true; do
    if [ ! -S "/runtime-mnt/vineyard/default/vineyard/vineyard-fuse/vineyard.sock" ]; then
        mount --bind /runtime-mnt/vineyard/default/vineyard /runtime-mnt/vineyard/default/vineyard/vineyard-fuse
    fi
    sleep 10
done