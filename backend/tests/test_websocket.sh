#!/bin/bash

cd $(dirname $0)

if [ ! -f ./websocat.x86_64-unknown-linux-musl ]; then
    echo "downloading websocat..."
    wget https://github.com/vi/websocat/releases/download/v1.12.0/websocat.x86_64-unknown-linux-musl
fi

# need to pass user secret in header
./websocat.x86_64-unknown-linux-musl wss://www.toymaker-ben.online/api/ws/1