#!/bin/bash

script_dir="$(cd "$(dirname "$0")" && pwd)"
cd ${script_dir}

echo "testing gRPC client ..."
go run .